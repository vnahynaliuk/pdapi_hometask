package handlers

import (
    "io"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/gorilla/mux"
    "github.com/vnahynaliuk/pdapi_hometask/middleware"
)

// setupTestRouter creates a router similar to the one in main.go,
// but without launching a real HTTP server on port 8080.
func setupTestRouter() *mux.Router {
    r := mux.NewRouter()

    // Attach middleware
    r.Use(middleware.LoggingMiddleware)
    r.Use(middleware.MetricsMiddleware)

    // Register handlers exactly as in main.go
    r.HandleFunc("/deals", GetDeals).Methods("GET")
    r.HandleFunc("/deals", AddDeal).Methods("POST")
    r.HandleFunc("/deals/{id}", UpdateDeal).Methods("PUT")

    return r
}

// TestGetDeals is a simple test that checks
// if GET /deals does not return a 404 error,
// and reads the response body.
func TestGetDeals(t *testing.T) {
    router := setupTestRouter()

    // Create a fake HTTP request
    req, err := http.NewRequest("GET", "/deals", nil)
    if err != nil {
        t.Fatal(err) // stops the test if the request creation failed
    }

    // httptest.NewRecorder is a “fake” response writer that captures the result
    rr := httptest.NewRecorder()

    // ServeHTTP passes the request to the router
    router.ServeHTTP(rr, req)

    // Check the status code
    if rr.Code == http.StatusNotFound {
        t.Errorf("Expected not to get 404. Got %d instead", rr.Code)
    }

    // Optionally, inspect the response body (for JSON, etc.)
    respBody, _ := io.ReadAll(rr.Body)
    t.Logf("GET /deals response body: %s", string(respBody))
}

// TestAddDeal is an example of a test for POST /deals.
// It checks that the route exists and does not return 404.
func TestAddDeal(t *testing.T) {
    router := setupTestRouter()

    // Construct a minimal JSON body for creating a "Deal"
    requestBody := `{"title":"Test Deal via TestAddDeal"}`
    req, err := http.NewRequest("POST", "/deals", strings.NewReader(requestBody))
    if err != nil {
        t.Fatal(err)
    }
    // Set the Content-Type to application/json
    req.Header.Set("Content-Type", "application/json")

    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    if rr.Code == http.StatusNotFound {
        t.Errorf("Expected not to get 404. Got 404 instead")
    }

    respBody, _ := io.ReadAll(rr.Body)
    t.Logf("POST /deals response body: %s", string(respBody))
}
