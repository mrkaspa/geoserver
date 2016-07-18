package ws

import (
	"encoding/json"
	"time"

	"github.com/mrkaspa/geoserver/models"
	"github.com/mrkaspa/geoserver/utils"
	"golang.org/x/net/websocket"
)

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn
	// name of the connection
	name string
	// actor reference
	actorRef *actor
	// Buffered channel of outbound messages.
	send chan []byte
	// channel for inbound messages
	receive chan []byte
	// channel to kill the connection
	poisonPill chan bool
}

func createConnection(name string, ws *websocket.Conn) *connection {
	return &connection{
		ws:         ws,
		name:       name,
		send:       make(chan []byte, 256),
		receive:    make(chan []byte, 256),
		poisonPill: make(chan bool),
	}
}

// readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump() {
	for {
		var message []byte

		if err := websocket.Message.Receive(c.ws, &message); err != nil {
			utils.Log.Infof("Can't receive")
			break
		}

		utils.Log.Infof("Message recived (%s)", string(message))
		c.receive <- message // sends to processMessages
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *connection) writePump() {
	for {
		select {
		case <-c.poisonPill:
			if len(c.send) == 0 {
				c.die()
				return
			}
			go c.dieLater()
		case message := <-c.send:
			// channel used to finish the connection when it's closed
			utils.Log.Infof("Writing (%s) to the connection: %s", string(message), c.name)
			if err := websocket.Message.Send(c.ws, message); err != nil {
				utils.Log.Infof("Can't send")
				break
			}
		}
	}
}

func (c *connection) processMessages() {
	for message := range c.receive {
		c.processMessage(message)
	}
}

// process each message
func (c *connection) processMessage(message []byte) {
	register := registerActor{
		name:     c.name,
		response: make(chan *actor),
	}
	SearcherVar.register <- &register
	actorRef := <-register.response
	actorRef.addConnection <- c
	//creates the postStroke
	strokeVar := new(models.Stroke)
	json.Unmarshal(message, strokeVar)
	strokeVar.UserID = actorRef.name
	actorRef.strokes <- strokeVar
	for response := range actorRef.responses {
		// expects all the responses from the actor until it dies
		c.send <- response
	}
}

func (c *connection) dieLater() {
	timer := time.NewTimer(100 * time.Millisecond)
	<-timer.C
	c.poisonPill <- true
}

func (c *connection) die() {
	utils.Log.Infof("Closing websocket: %s", c.name)
	c.ws.Close()
}
