package ws

import (
	"net/http"

	"github.com/mrkaspa/geoserver/utils"

	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
)

func Adapter(w http.ResponseWriter, req *http.Request) {
	s := websocket.Server{Handler: websocket.Handler(handler)}
	s.ServeHTTP(w, req)
}

// Handler for ws
func handler(ws *websocket.Conn) {
	r := ws.Request()
	vars := mux.Vars(r)
	username := vars["username"]
	utils.Log.Infof("Creating connection: %s", username)
	c := createConnection(username, ws)
	go c.writePump()
	go c.processMessages()
	c.readPump()
}
