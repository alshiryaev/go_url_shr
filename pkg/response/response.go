package response

import (
	"encoding/json"
	"net/http"
)

func Json(w http.ResponseWriter, payload any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}
