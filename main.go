// Package main
//
// @title Pipedrive Deals API
// @version 1.0
// @description A simple proxy to Pipedrive API for managing deals.

// @host localhost:8080
// @BasePath /
// @schemes http
//
// @securityDefinitions.apikey ApiKeyAuth
// @in query
// @name api_token
package main

import (
    "bytes"
    _"encoding/json"
    "io"
    "log"
    "net/http"
    "os"
    _"strconv"
    "time"

    "github.com/gorilla/mux"
    httpSwagger "github.com/swaggo/http-swagger"
    _ "github.com/vnahynaliuk/pdapi_hometask/docs" // replace with your actual module path
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// Read API token and company domain from environment variables.
var (
    apiToken      = os.Getenv("PIPEDRIVE_API_TOKEN")
    companyDomain = os.Getenv("PIPEDRIVE_COMPANY_DOMAIN")
)

// Define Prometheus metrics.
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

func main() {
    // Register metrics only once.
    prometheus.MustRegister(httpRequestsTotal, httpRequestDuration)

    r := mux.NewRouter()
    r.Use(loggingMiddleware)
    r.Use(metricsMiddleware)

    // GET and POST endpoints for deals.
    r.HandleFunc("/deals", getDeals).Methods("GET")
    r.HandleFunc("/deals", addDeal).Methods("POST")
    // PUT endpoint now expects the deal ID as a path parameter.
    r.HandleFunc("/deals/{id}", updateDeal).Methods("PUT")

    // Expose Prometheus metrics at /metrics.
    r.Handle("/metrics", promhttp.Handler())
    // Swagger UI endpoint.
    r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

    log.Println("Server started on port 8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}

// loggingMiddleware logs the request method, URI, and remote address.
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
        next.ServeHTTP(w, r)
    })
}

// metricsMiddleware collects metrics for each request.
func metricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues(r.Method, r.RequestURI))
        defer timer.ObserveDuration()
        httpRequestsTotal.WithLabelValues(r.Method, r.RequestURI).Inc()
        next.ServeHTTP(w, r)
    })
}

// forwardRequest creates and sends an HTTP request to the Pipedrive API.
func forwardRequest(method, url string, body io.Reader) (*http.Response, error) {
    req, err := http.NewRequest(method, url, body)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    client := &http.Client{Timeout: 10 * time.Second}
    return client.Do(req)
}

// CreateDeal represents the request body for creating a deal.
// Fields and examples are based on Pipedrive API documentation.
type AddDeal struct {
    Title             string   `json:"title" example:"Test Deal"`                              // REQUIRED
    Value             string   `json:"value" example:"1000"`
    Label             []int    `json:"label" example:"[1,2,3]"`
    Currency          string   `json:"currency" example:"USD"`
    UserID            int      `json:"user_id" example:"123"`
    PersonID          int      `json:"person_id" example:"456"`
    OrgID             int      `json:"org_id" example:"789"`
    PipelineID        int      `json:"pipeline_id" example:"10"`
    StageID           int      `json:"stage_id" example:"20"`
    Status            string   `json:"status" example:"open" enums:"open,won,lost,deleted"`
    OriginID          string   `json:"origin_id" example:"integration_xyz"`
    Channel           int      `json:"channel" example:"1"`
    ChannelID         string   `json:"channel_id" example:"ch_123"`
    AddTime           string   `json:"add_time" example:"2023-08-21 12:34:56"`
    WonTime           string   `json:"won_time" example:"2023-08-22 13:00:00"`
    LostTime          string   `json:"lost_time" example:"2023-08-22 14:00:00"`
    CloseTime         string   `json:"close_time" example:"2023-08-23 15:00:00"`
    ExpectedCloseDate string   `json:"expected_close_date" example:"2023-08-30"`
    Probability       float64  `json:"probability" example:"75.5"`
    LostReason        string   `json:"lost_reason" example:"Price too high"`
    VisibleTo         string   `json:"visible_to" example:"1" enums:"1,3,5,7"`
}

// UpdateDeal represents the request body for updating a deal.
// All fields are optional.
type UpdateDeal struct {
    Title             *string  `json:"title,omitempty" example:"Updated Deal Title"`
    Value             *string  `json:"value,omitempty" example:"1500"`
    Label             []int    `json:"label,omitempty" example:"[1,2]"`
    Currency          *string  `json:"currency,omitempty" example:"USD"`
    UserID            *int     `json:"user_id,omitempty" example:"123"`
    PersonID          *int     `json:"person_id,omitempty" example:"456"`
    OrgID             *int     `json:"org_id,omitempty" example:"789"`
    PipelineID        *int     `json:"pipeline_id,omitempty" example:"10"`
    StageID           *int     `json:"stage_id,omitempty" example:"20"`
    Status            *string  `json:"status,omitempty" example:"open" enums:"open,won,lost,deleted"`
    Channel           *int     `json:"channel,omitempty" example:"1"`
    ChannelID         *string  `json:"channel_id,omitempty" example:"ch_123"`
    WonTime           *string  `json:"won_time,omitempty" example:"2023-08-22 13:00:00"`
    LostTime          *string  `json:"lost_time,omitempty" example:"2023-08-22 14:00:00"`
    CloseTime         *string  `json:"close_time,omitempty" example:"2023-08-23 15:00:00"`
    ExpectedCloseDate *string  `json:"expected_close_date,omitempty" example:"2023-08-30"`
    Probability       *float64 `json:"probability,omitempty" example:"75.5"`
    LostReason        *string  `json:"lost_reason,omitempty" example:"Price too high"`
    VisibleTo         *string  `json:"visible_to,omitempty" example:"1" enums:"1,3,5,7"`
}

// getDeals godoc
// @Summary Retrieve all deals
// @Description Retrieves all deals from the Pipedrive API with optional query parameters.
// @Produce json
// @Security ApiKeyAuth
// (Query parameters for GET are defined separately. See previous documentation.)
// @Success 200 {object} map[string]interface{}
// @Router /deals [get]
func getDeals(w http.ResponseWriter, r *http.Request) {
    baseURL := "https://" + companyDomain + ".pipedrive.com/api/v1/deals?api_token=" + apiToken

    // Append any query parameters from the client request.
    q := r.URL.Query()
    for key, values := range q {
        for _, value := range values {
            baseURL += "&" + key + "=" + value
        }
    }

    resp, err := forwardRequest("GET", baseURL, nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    w.Header().Set("Content-Type", "application/json")
    io.Copy(w, resp.Body)
}

// addDeal godoc
// @Summary Create a new deal
// @Description Creates a new deal via the Pipedrive API.
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param deal body CreateDeal true "Deal data"
// @Success 200 {object} map[string]interface{}
// @Router /deals [post]
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

// updateDeal godoc
// @Summary Update an existing deal
// @Description Updates an existing deal via the Pipedrive API.
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "The ID of the deal"
// @Param deal body UpdateDeal true "Deal update data"
// @Success 200 {object} map[string]interface{}
// @Router /deals/{id} [put]
func updateDeal(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    bodyBytes, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    url := "https://" + companyDomain + ".pipedrive.com/api/v1/deals/" + id + "?api_token=" + apiToken
    resp, err := forwardRequest("PUT", url, bytes.NewBuffer(bodyBytes))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    w.Header().Set("Content-Type", "application/json")
    io.Copy(w, resp.Body)
}
