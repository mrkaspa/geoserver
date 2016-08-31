package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mrkaspa/geoserver/models"
	"github.com/urfave/negroni"
)

type fakeRouter struct {
	mux *mux.Router
}

func NewRouter() http.Handler {
	n := negroni.Classic()
	router := &fakeRouter{mux: mux.NewRouter()}
	controller := &controller{persistorCreator: models.NewPersistor}
	router.mux.HandleFunc("/near", controller.nearHandler).Methods("POST")
	router.mux.HandleFunc("/store", controller.storeHandler).Methods("POST")
	router.mux.HandleFunc("/recent/{id}", controller.recentStrokes).Methods("GET")
	n.Use(router)
	return n
}

func (router *fakeRouter) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	router.mux.ServeHTTP(w, r)
	next(w, r)
}
