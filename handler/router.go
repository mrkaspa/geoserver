package handler

import (
	"github.com/gorilla/mux"
	"github.com/mrkaspa/geoserver/ws"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/ws/{username}", ws.Adapter)
	router.HandleFunc("/near/{username}", nearHandler).Methods("GET")
	router.HandleFunc("/store", storeHandler).Methods("POST")
	return router
}
