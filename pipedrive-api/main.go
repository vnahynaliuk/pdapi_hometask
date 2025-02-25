package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var apiToken = os.Getenv("PIPEDRIVE_API_TOKEN")
var companyDomain = os.Getenv("PIPEDRIVE_COMPANY_DOMAIN")

func getDeals(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://" + companyDomain + ".pipedrive.com/api/v1/deals?api_token=" + apiToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func addDeal(w http.ResponseWriter, r *http.Request) {
	// Реалізація додавання угоди
}

func updateDeal(w http.ResponseWriter, r *http.Request) {
	// Реалізація оновлення угоди
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/deals", getDeals).Methods("GET")
	r.HandleFunc("/deals", addDeal).Methods("POST")
	r.HandleFunc("/deals", updateDeal).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8080", r))
}
