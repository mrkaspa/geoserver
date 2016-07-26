package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mrkaspa/geoserver/models"
	"github.com/mrkaspa/geoserver/utils"
)

type controller struct {
	persistor models.Persistance
}

func (c *controller) nearHandler(w http.ResponseWriter, r *http.Request) {
	strokeNear := new(models.StrokeNear)
	utils.Log.Infof("/near: %v", strokeNear)
	err := json.NewDecoder(r.Body).Decode(strokeNear)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	c.persistor.PersistAndFind() <- strokeNear
	saved := <-c.persistor.Saved()
	nearUsers := <-c.persistor.UsersFound()
	if !saved {
		utils.Log.Infof("Response %d", http.StatusConflict)
		w.WriteHeader(http.StatusConflict)
		return
	}
	sendOkJSON(w, nearUsers)
}

func (c *controller) storeHandler(w http.ResponseWriter, r *http.Request) {
	stroke := new(models.Stroke)
	err := json.NewDecoder(r.Body).Decode(stroke)
	if err != nil {
		utils.Log.Infof("Response %d", http.StatusNotAcceptable)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	c.persistor.Persist() <- stroke
	saved := <-c.persistor.Saved()
	if !saved {
		utils.Log.Infof("Response %d", http.StatusConflict)
		w.WriteHeader(http.StatusConflict)
		return
	}
	utils.Log.Infof("Response %d", http.StatusOK)
	w.WriteHeader(http.StatusOK)
}
