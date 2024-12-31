package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// SMS request struct
type SMSRequest struct {
	Number  string `json:"number"`
	Message string `json:"message"`
	Name    string `json:"name,omitempty"`
}

func sendSMS(w http.ResponseWriter, r *http.Request) {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	clientID := os.Getenv("SMS_CLIENT_ID")
	clientSecret := os.Getenv("SMS_CLIENT_SECRET")
	gatewayURL := os.Getenv("SMS_GATEWAY_URL")

	// Parse the request body
	var smsReq SMSRequest
	err = json.NewDecoder(r.Body).Decode(&smsReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Prepare the data to send to the gateway
	data := map[string]interface{}{
		"message": smsReq.Message,
		"to":      smsReq.Number,
		"sender":  smsReq.Name,
	}

	// Create JSON from the data
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Failed to create JSON", http.StatusInternalServerError)
		return
	}

	// Make the API call to the SMS gateway
	req, err := http.NewRequest("POST", gatewayURL, bytes.NewBuffer(jsonData))
	if err != nil {
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
		http.Error(w, "Error calling SMS gateway", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Return the response from the SMS gateway
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("SMS sent successfully"))
}

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Set up router
	r := mux.NewRouter()
	r.HandleFunc("/send-sms", sendSMS).Methods("POST")

	port := os.Getenv("APP_PORT")
	if port == "" {
		log.Fatal("APP_PORT environment variable is not set")
	}
	// Start server
	log.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
