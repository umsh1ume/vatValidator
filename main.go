package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"vatValidator/services"
)

type VatValidationRequestData struct {
	VatNumber string `json:"vat_number"`
}

type VatValidationResponseData struct {
	VatNumber string `json:"vat_number"`
	IsValid   bool   `json:"is_valid"`
	Message   string `json:"message"`
}

func NewVatValidationResponse(vatNumber string, isValid bool, message string) VatValidationResponseData {
	return VatValidationResponseData{
		VatNumber: vatNumber,
		IsValid:   isValid,
		Message:   message,
	}
}

func main() {
	http.HandleFunc("/validate_vat", ValidateVat)
	if err := http.ListenAndServe(":8001", nil); err != nil {
		log.Fatal(err)
	}
}

// to handle CORS and preflight request
func setupResponseHeaders(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// controller function
func ValidateVat(w http.ResponseWriter, r *http.Request) {

	setupResponseHeaders(&w)
	if r.Method == "OPTIONS" {
		return
	}

	// Double check it's a post request being made
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}

	decoder := json.NewDecoder(r.Body)

	var requestData VatValidationRequestData
	err := decoder.Decode(&requestData)
	if err != nil || requestData.VatNumber == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid request")
		return
	}

	isValid, errMsg, validateErr := services.ValidateGermanVat(requestData.VatNumber)

	if validateErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, errMsg)
		return
	}

	response := NewVatValidationResponse(requestData.VatNumber, isValid, errMsg)
	w.WriteHeader(http.StatusOK)
	responseBytes, _ := json.Marshal(response)
	fmt.Fprintf(w, string(responseBytes))

	return
}
