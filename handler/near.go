package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mrkaspa/geoserver/models"
	"github.com/mrkaspa/geoserver/utils"
)

func nearHandler(w http.ResponseWriter, r *http.Request) {
	stroke := new(models.Stroke)
	utils.Log.Infof("/near: %v", stroke)
	err := json.NewDecoder(r.Body).Decode(stroke)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	persistor := models.NewPersistor()
	persistor.PersistAndFind <- stroke
	saved := <-persistor.Saved
	nearUsers := <-persistor.UsersFound
	if !saved {
		utils.Log.Infof("Response %d", http.StatusConflict)
		w.WriteHeader(http.StatusConflict)
		return
	}
	sendOkJSON(w, nearUsers)
}
