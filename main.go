package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	_ "sms_service/docs"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	swagger "github.com/swaggo/http-swagger"
)

// SMSRequest defines the structure for the SMS request
// @Description SMS request structure with number, message, and sender
// SMSRequest represents the structure of the SMS request payload
type SMSRequest struct {
	Number    interface{} `json:"number"`              // @param number query Use the formatted "to" field (either single or array) "Phone number. The number(s) that will receive the message"
	Message   string      `json:"message"`             // @param message query string true "The message to be sent"
	Sender    string      `json:"sender"`              // @param sender query string true "Sender. The number or name of the sender. A number can't be longer than 14 characters. A name can't be longer than 11 characters and can't contain special characters or spaces"
	Date      string      `json:"date,omitempty"`      // @param date query string false "Date and time the message will be sent in yyyy-MM-dd HH:mm format. If not provided, the message will be sent as soon as possible"
	Reference string      `json:"reference,omitempty"` // @param reference query string false "Custom reference. A string of max. 255 characters"
	Test      bool        `json:"test,omitempty"`      // @param test query bool false "If true, the system will check all parameters but will not send an SMS message (no credits/balance used)"
	Subid     string      `json:"subid,omitempty"`     // @param subid query string false "ID of a subaccount. If specified, the message will be sent from the subaccount"
}

var errorLog *log.Logger

// sendSMS handles sending SMS to a number
// @Summary Send SMS
// @Description Sends an SMS message to a specified phone number
// @Accept  json
// @Produce  json
// @Param sms body SMSRequest true "Send SMS request"
// @Success 200 {string} string "SMS sent successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Router /send-sms [post]
func sendSMS(w http.ResponseWriter, r *http.Request) {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		errorLog.Printf("Error loading .env file: %v", err)
		http.Error(w, "Error loading .env file", http.StatusInternalServerError)
		return
	}

	clientID := os.Getenv("SMS_CLIENT_ID")
	clientSecret := os.Getenv("SMS_CLIENT_SECRET")
	gatewayURL := os.Getenv("SMS_GATEWAY_URL")

	// Parse the request body
	var smsReq SMSRequest
	err = json.NewDecoder(r.Body).Decode(&smsReq)
	if err != nil {
		errorLog.Printf("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Handle the "to" field: if it's a single string, convert it to an array
	var toNumbers interface{}
	switch v := smsReq.Number.(type) {
	case string:
		toNumbers = []string{v} // Convert single number to an array
	case []interface{}:
		toNumbers = v // Already an array, use it as is
	default:
		errorLog.Printf("Invalid number format: %v", smsReq.Number)
		http.Error(w, "Invalid number format", http.StatusBadRequest)
		return
	}

	// Prepare the data to send to the gateway
	data := map[string]interface{}{
		"message":   smsReq.Message,
		"to":        toNumbers,
		"sender":    smsReq.Sender,
		"date":      smsReq.Date,      // Use the Date field from the request if available
		"reference": smsReq.Reference, // Use the Reference field from the request if available
		"test":      smsReq.Test,      // Use the Test field from the request (boolean)
	}

	// Only add subid if it's provided
	if smsReq.Subid != "" {
		data["subid"] = smsReq.Subid // Add the subid field if it's not empty
	}

	// Create JSON from the data
	jsonData, err := json.Marshal(data)
	if err != nil {
		errorLog.Printf("Failed to create JSON: %v", err)
		http.Error(w, "Failed to create JSON", http.StatusInternalServerError)
		return
	}

	// Make the API call to the SMS gateway
	req, err := http.NewRequest("POST", gatewayURL, bytes.NewBuffer(jsonData))
	if err != nil {
		errorLog.Printf("Error creating request: %v", err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Client-Id", clientID)
	req.Header.Set("X-Client-Secret", clientSecret)

	// Make the request to the SMS gateway
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		errorLog.Printf("Error calling SMS gateway: %v", err)
		http.Error(w, "Error calling SMS gateway", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read and log the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		errorLog.Printf("Error reading response: %v", err)
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		return
	}
	log.Printf("Response from SMS Gateway: %s", string(responseBody))

	// Return the response from the SMS gateway
	if resp.StatusCode != http.StatusOK {
		errorLog.Printf("Failed to send message: %s", string(responseBody))
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("SMS sent successfully"))
}

// healthCheck returns a simple health check status
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("API is up and running"))
}

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Set up log file for error logging
	logFile, err := os.OpenFile("sms_service_app_error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	defer logFile.Close()

	// Create a new logger that writes to the log file
	errorLog = log.New(logFile, "ERROR: ", log.LstdFlags|log.Lshortfile)

	// Set up router
	r := mux.NewRouter()
	r.HandleFunc("/send-sms", sendSMS).Methods("POST")

	r.HandleFunc("/health", healthCheck).Methods("GET")

	// Swagger UI route
	r.PathPrefix("/swagger/").Handler(swagger.WrapHandler)

	port := os.Getenv("APP_PORT")
	if port == "" {
		log.Fatal("APP_PORT environment variable is not set")
	}
	// Start server
	log.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
