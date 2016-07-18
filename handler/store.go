package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mrkaspa/geoserver/models"
)

func storeHandler(w http.ResponseWriter, r *http.Request) {
	stroke := new(models.Stroke)
	err := json.NewDecoder(r.Body).Decode(stroke)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
		return
	}

}
