package ws

import (
	"time"

	"math"

	"github.com/mrkaspa/geoserver/models"
	"github.com/mrkaspa/geoserver/utils"
)

const defaultLifeTime = 1

type actorStatus int

type broadcastActor struct {
	actor     *actor
	recursive bool
}

const (
	_                 = iota
	alive actorStatus = 1
	dead
)

// Represents an user in the system that is doing a cheers
type actor struct {
	name             string
	status           actorStatus
	info             []byte
	timer            *time.Timer
	lifeTime         time.Time
	connections      []*connection
	persistor        models.Persistance
	matchedActors    map[*actor]bool
	sentActors       map[*actor]bool
	addConnection    chan *connection
	removeConnection chan *connection
	strokesNear      chan *models.StrokeNear
	responses        chan []byte
	broadcast        chan *broadcastActor
	ping             chan *actor
	pong             chan *actor
	poisonPill       chan bool
}

// inits an actor
func newActor(name string) *actor {
	utils.Log.Infof("Creating actor: %s", name)
	actor := &actor{
		name:             name,
		status:           alive,
		connections:      []*connection{},
		lifeTime:         time.Now(),
		persistor:        models.NewPersistor(),
		addConnection:    make(chan *connection),
		removeConnection: make(chan *connection),
		matchedActors:    make(map[*actor]bool),
		sentActors:       make(map[*actor]bool),
		strokesNear:      make(chan *models.StrokeNear),
		responses:        make(chan []byte),
		broadcast:        make(chan *broadcastActor, 256),
		ping:             make(chan *actor, 256),
		pong:             make(chan *actor, 256),
		poisonPill:       make(chan bool, 1),
	}
	go actor.startTimer(defaultLifeTime)
	go actor.run()
	return actor
}

// run the actor that represents the user
func (a *actor) run() {
	defer func() {
		if r := recover(); r != nil {
			utils.Log.Infof("Recovered in %s.run()", a.name)
		}
	}()

	utils.Log.Infof("Running actor: %s", a.name)
	for {
		select {
		case <-a.poisonPill:
			a.die()
			return
		case conn := <-a.addConnection:
			utils.Log.Infof("adding connection: %s to: %s - len before add: %d", conn.name, a.name, len(a.connections))
			a.connections = append(a.connections, conn)
		case conn := <-a.removeConnection:
			a.removeConnectionBy(conn)
		case strokeVar := <-a.strokesNear:
			a.persist(strokeVar)
		case users := <-a.persistor.UsersFound():
			a.receivingNearUsers(users)
		case actorPing := <-a.ping:
			a.receivingPing(actorPing)
		case actorPong := <-a.pong:
			a.receivingPong(actorPong)
		case actorToSend := <-a.broadcast:
			a.sendBroadcast(actorToSend)
		}
	}
}

func (a *actor) receivingNearUsers(users []models.Stroke) {
	for _, u := range users {
		searchActorVar := searchActor{
			name:     u.UserID,
			response: make(chan *actor),
		}
		SearcherVar.search <- &searchActorVar
		actorRef := <-searchActorVar.response
		a.matchedActors[actorRef] = false
		if actorRef != nil {
			actorRef.ping <- a
		}
	}
}

func (a *actor) receivingPing(actorPing *actor) {
	a.broadcast <- &broadcastActor{actorPing, true}
	_, ok := a.matchedActors[actorPing]
	utils.Log.Infof("%s received a PING from %s, ok: %t", a.name, actorPing.name, ok)
	if !ok {
		utils.Log.Infof("%s entered a PING from %s", a.name, actorPing.name)
		a.matchedActors[actorPing] = false
		diffTime := actorPing.lifeTime.Sub(a.lifeTime)
		moreTime := int(math.Ceil(diffTime.Seconds()))
		a.dieLater(moreTime)
		actorPing.pong <- a
	}
}

func (a *actor) receivingPong(actorPong *actor) {
	utils.Log.Infof("%s received a PONG from %s", a.name, actorPong.name)
	a.matchedActors[actorPong] = true
	a.broadcast <- &broadcastActor{actorPong, true}
}

func (a *actor) sendBroadcast(actorToSend *broadcastActor) {
	actorFound := actorToSend.actor
	a.broadcastActor(actorFound)
	utils.Log.Infof("sendBroadcast a: %s, actorFound: %s, recursive: %t", a.name, actorFound.name, actorToSend.recursive)
	if actorToSend.recursive {
		for actorPonged, ponged := range a.matchedActors {
			if ponged && actorPonged != actorFound {
				utils.Log.Infof("sendBroadcast a: %s, matchedActors: %s", a.name, actorPonged.name)
				actorPonged.broadcast <- &broadcastActor{actorFound, false}
			}
		}
		for actorPonged, ponged := range actorFound.matchedActors {
			if ponged && actorPonged != a {
				utils.Log.Infof("sendBroadcast actorFound: %s, matchedActors: %s", actorFound.name, actorPonged.name)
				a.broadcastActor(actorPonged)
				actorPonged.broadcast <- &broadcastActor{a, false}
			}
		}
	}
}

// sends the persist message
func (a *actor) persist(strokeNear *models.StrokeNear) {
	a.info = []byte(strokeNear.Stroke.Info)
	a.persistor.PersistAndFind() <- strokeNear
	if <-a.persistor.Saved() {
		go a.dieLater(strokeNear.TimeRange)
	}
}

// sends data to all connectios
func (a *actor) broadcastActor(actorMatched *actor) {
	ok := a.sentActors[actorMatched]
	sent := false
	if !ok {
		sent = true
		a.sentActors[actorMatched] = true
		for _, conn := range a.connections {
			utils.Log.Infof("SendingBC to the connection %s the info %s of the actor %s", a.name, string(actorMatched.info), actorMatched.name)
			conn.send <- actorMatched.info
		}
	}
	utils.Log.Infof("Sent? %t, %s broadcasting %s, found: %t, on connections: %d", sent, a.name, actorMatched.name, ok, len(a.connections))
}

// removes a connection
func (a *actor) removeConnectionBy(conn *connection) {
	for i, c := range a.connections {
		if c == conn {
			a.connections = append(a.connections[:i], a.connections[i+1:]...)
			return
		}
	}
}

// timeout to kill the actor started after the actor is registered
func (a *actor) startTimer(seconds int) {
	a.lifeTime = a.lifeTime.Add(time.Duration(seconds) * time.Second)
	a.timer = time.NewTimer(time.Duration(seconds) * time.Second)
	<-a.timer.C
	a.poisonPill <- true
}

// postpones death
func (a *actor) dieLater(seconds int) {
	if seconds > 0 {
		a.timer.Stop()
		a.timer = nil
		utils.Log.Infof("Actor %s die later %d", a.name, seconds)
		go a.startTimer(seconds)
	}
}

// finish an actor
func (a *actor) die() {
	// kills the referenced actors
	if a.status == alive {
		utils.Log.Infof("Actor dying: %s -- with %d connections -- should die %v is dying %v", a.name, len(a.connections), a.lifeTime, time.Now())

		a.status = dead

		// closes all the actor connections
		for _, conn := range a.connections {
			utils.Log.Infof("Closing connection: %s", conn.name)
			conn.poisonPill <- true
			a.removeConnection <- conn
		}
		SearcherVar.unregister <- a.name
	}
}
