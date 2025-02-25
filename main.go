package main

import (
    "bytes"
    "encoding/json"
    "io"
    "log"
    "net/http"
    "os"
    "strconv"
    "time"

    "github.com/gorilla/mux"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// Read API token and company domain from environment variables
var (
    apiToken      = os.Getenv("PIPEDRIVE_API_TOKEN")
    companyDomain = os.Getenv("PIPEDRIVE_COMPANY_DOMAIN")
)

// Set up Prometheus metrics
var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests.",
        },
        []string{"method", "endpoint"},
    )
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "Duration of HTTP requests in seconds.",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )
)

func init() {
    prometheus.MustRegister(httpRequestsTotal, httpRequestDuration)
}

// Logging middleware that logs the request method, URI, and remote address.
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
        next.ServeHTTP(w, r)
    })
}

// Metrics middleware collects metrics for each request.
func metricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues(r.Method, r.RequestURI))
        defer timer.ObserveDuration()
        httpRequestsTotal.WithLabelValues(r.Method, r.RequestURI).Inc()
        next.ServeHTTP(w, r)
    })
}

// forwardRequest creates and sends an HTTP request with the specified method, URL, and body.
func forwardRequest(method, url string, body io.Reader) (*http.Response, error) {
    req, err := http.NewRequest(method, url, body)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    client := &http.Client{Timeout: 10 * time.Second}
    return client.Do(req)
}

// GET /deals - Retrieves all deals from the Pipedrive API.
func getDeals(w http.ResponseWriter, r *http.Request) {
    url := "https://" + companyDomain + ".pipedrive.com/api/v1/deals?api_token=" + apiToken
    resp, err := forwardRequest("GET", url, nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    w.Header().Set("Content-Type", "application/json")
    io.Copy(w, resp.Body)
}

// POST /deals - Creates a new deal via the Pipedrive API.
func addDeal(w http.ResponseWriter, r *http.Request) {
    url := "https://" + companyDomain + ".pipedrive.com/api/v1/deals?api_token=" + apiToken
    bodyBytes, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    resp, err := forwardRequest("POST", url, bytes.NewBuffer(bodyBytes))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    w.Header().Set("Content-Type", "application/json")
    io.Copy(w, resp.Body)
}

// PUT /deals - Updates an existing deal. Expects an "id" field in the JSON body.
func updateDeal(w http.ResponseWriter, r *http.Request) {
    var data map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    idFloat, ok := data["id"].(float64)
    if !ok {
        http.Error(w, "Missing or invalid 'id' field", http.StatusBadRequest)
        return
    }
    id := strconv.Itoa(int(idFloat))
    // Remove the "id" field if the API does not expect it in the JSON body.
    delete(data, "id")
    jsonData, err := json.Marshal(data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    url := "https://" + companyDomain + ".pipedrive.com/api/v1/deals/" + id + "?api_token=" + apiToken
    resp, err := forwardRequest("PUT", url, bytes.NewBuffer(jsonData))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    w.Header().Set("Content-Type", "application/json")
    io.Copy(w, resp.Body)
}

func main() {
    r := mux.NewRouter()

    // Add logging and metrics middleware.
    r.Use(loggingMiddleware)
    r.Use(metricsMiddleware)

    // Register API endpoints.
    r.HandleFunc("/deals", getDeals).Methods("GET")
    r.HandleFunc("/deals", addDeal).Methods("POST")
    r.HandleFunc("/deals", updateDeal).Methods("PUT")
    // Expose the Prometheus metrics endpoint.
    r.Handle("/metrics", promhttp.Handler())

    log.Println("Server started on port 8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}