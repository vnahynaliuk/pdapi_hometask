package handlers

import (
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/vnahynaliuk/pdapi_hometask/utils"
)

// Get environment variables for API token and company domain.
var (
	apiToken      = os.Getenv("PIPEDRIVE_API_TOKEN")
	companyDomain = os.Getenv("PIPEDRIVE_COMPANY_DOMAIN")
)

// GetDeals handles GET requests to retrieve deals with optional query parameters.
//
// @Summary Retrieve all deals
// @Description Retrieves all deals from the Pipedrive API with optional query parameters.
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /deals [get]
func GetDeals(w http.ResponseWriter, r *http.Request) {
	baseURL := "https://" + companyDomain + ".pipedrive.com/api/v1/deals?api_token=" + apiToken

	// Append any query parameters from the client request.
	q := r.URL.Query()
	for key, values := range q {
		for _, value := range values {
			baseURL += "&" + key + "=" + value
		}
	}

	resp, err := utils.ForwardRequest("GET", baseURL, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}

// AddDeal handles POST requests to create a new deal.
//
// @Summary Create a new deal
// @Description Creates a new deal via the Pipedrive API.
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param deal body models.CreateDeal true "Deal data"
// @Success 200 {object} map[string]interface{}
// @Router /deals [post]
func AddDeal(w http.ResponseWriter, r *http.Request) {
	url := "https://" + companyDomain + ".pipedrive.com/api/v1/deals?api_token=" + apiToken
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := utils.ForwardRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}

// UpdateDeal handles PUT requests to update an existing deal.
//
// @Summary Update an existing deal
// @Description Updates an existing deal via the Pipedrive API.
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "The ID of the deal"
// @Param deal body models.UpdateDeal true "Deal update data"
// @Success 200 {object} map[string]interface{}
// @Router /deals/{id} [put]
func UpdateDeal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	url := "https://" + companyDomain + ".pipedrive.com/api/v1/deals/" + id + "?api_token=" + apiToken
	resp, err := utils.ForwardRequest("PUT", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}
