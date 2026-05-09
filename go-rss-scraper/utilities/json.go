package utilities

import (
	"encoding/json"
	"log"
	"net/http"
)

// Error Response
//
// {
// "error": "Something went wrong"
// }
type errrResponse struct {
	Error string `json:"error"`
}

// Responsd with JSON Utility
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

// Respond with Error
func RespondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5xx error: ", msg)
	}
	RespondWithJSON(w, code, errrResponse{
		Error: msg,
	})
}
