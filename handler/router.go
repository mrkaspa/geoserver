package handler

import (
	"github.com/gorilla/mux"
	"github.com/mrkaspa/geoserver/ws"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/ws/{username}", ws.ServeWS)
	router.HandleFunc("/near", nearHandler).Methods("GET")
	return router
}
