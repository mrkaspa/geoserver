package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mrkaspa/geoserver/models"
	"github.com/mrkaspa/geoserver/utils"
)

func nearHandler(w http.ResponseWriter, r *http.Request) {
	strokeNear := new(models.StrokeNear)
	utils.Log.Infof("/near: %v", strokeNear)
	err := json.NewDecoder(r.Body).Decode(strokeNear)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	persistor := models.NewPersistor()
	persistor.PersistAndFind <- strokeNear
	saved := <-persistor.Saved
	nearUsers := <-persistor.UsersFound
	if !saved {
		utils.Log.Infof("Response %d", http.StatusConflict)
		w.WriteHeader(http.StatusConflict)
		return
	}
	sendOkJSON(w, nearUsers)
}
