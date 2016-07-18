package handler

import (
	"encoding/json"
	"net/http"
)

func sendJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
