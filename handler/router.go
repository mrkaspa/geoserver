package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mrkaspa/geoserver/ws"
	"github.com/urfave/negroni"
)

type fakeRouter struct {
	mux *mux.Router
}

func NewRouter() http.Handler {
	n := negroni.Classic()
	router := &fakeRouter{mux: mux.NewRouter()}
	router.mux.HandleFunc("/ws/{username}", ws.Adapter)
	router.mux.HandleFunc("/near", nearHandler).Methods("POST")
	router.mux.HandleFunc("/store", storeHandler).Methods("POST")
	n.Use(router)
	return n
}

func (router *fakeRouter) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	router.mux.ServeHTTP(w, r)
	next(w, r)
}
