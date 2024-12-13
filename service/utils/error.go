package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	structs "github.com/edmii/WasaProject/service/models"
)

// SendErrorResponse writes a structured error response to the HTTP response writer
func SendErrorResponse(w http.ResponseWriter, message string, details []string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := structs.ErrorResponse{
		Status:  fmt.Sprintf("error %d", statusCode),
		Message: message,
		Details: details,
	}

	_ = json.NewEncoder(w).Encode(response)
}
