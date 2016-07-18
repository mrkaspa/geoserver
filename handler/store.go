package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mrkaspa/geoserver/models"
	"github.com/mrkaspa/geoserver/utils"
)

func storeHandler(w http.ResponseWriter, r *http.Request) {
	stroke := new(models.Stroke)
	err := json.NewDecoder(r.Body).Decode(stroke)
	if err != nil {
		utils.Log.Infof("Response %d", http.StatusNotAcceptable)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	persistor := models.NewPersistor()
	persistor.Persist <- stroke
	saved := <-persistor.Saved
	if !saved {
		utils.Log.Infof("Response %d", http.StatusConflict)
		w.WriteHeader(http.StatusConflict)
		return
	}
	utils.Log.Infof("Response %d", http.StatusOK)
	w.WriteHeader(http.StatusOK)
}
