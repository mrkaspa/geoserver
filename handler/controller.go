package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mrkaspa/geoserver/models"
	"github.com/mrkaspa/geoserver/utils"
)

type controller struct {
	persistorCreator func() models.Persistance
}

func (c controller) nearHandler(w http.ResponseWriter, r *http.Request) {
	strokeNear := models.StrokeNear{}
	err := json.NewDecoder(r.Body).Decode(&strokeNear)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	savedResponse := make(chan bool)
	usersResponse := make(chan []models.Stroke)
	persistor := c.persistorCreator()
	persistor.PersistAndFind() <- models.PersistAndFindWithResponse{
		StrokeNear:    strokeNear,
		SavedResponse: savedResponse,
		UsersResponse: usersResponse,
	}
	saved := <-savedResponse
	nearUsers := <-usersResponse

	if !saved {
		utils.Log.Infof("Response %d", http.StatusConflict)
		w.WriteHeader(http.StatusConflict)
		return
	}
	sendOkJSON(w, nearUsers)
}

func (c controller) storeHandler(w http.ResponseWriter, r *http.Request) {
	stroke := models.Stroke{}
	err := json.NewDecoder(r.Body).Decode(&stroke)
	if err != nil {
		utils.Log.Infof("Response %d", http.StatusNotAcceptable)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	response := make(chan bool)
	persistor := c.persistorCreator()
	persistor.Persist() <- models.PersistWithResponse{
		Stroke:   stroke,
		Response: response,
	}

	saved := <-response
	if !saved {
		utils.Log.Infof("Response %d", http.StatusConflict)
		w.WriteHeader(http.StatusConflict)
		return
	}
	utils.Log.Infof("Response %d", http.StatusOK)
	w.WriteHeader(http.StatusOK)
}

func (c controller) recentStrokes(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["id"]
	persistor := c.persistorCreator()
	strokesResponse := make(chan []models.Stroke)
	persistor.FindStrokes() <- models.FindStrokesWithResponse{
		Username: username,
		Response: strokesResponse,
	}
	history := <-strokesResponse
	sendOkJSON(w, history)
}
