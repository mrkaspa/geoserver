package handler

type searchActorWithResponse struct {
	name     string
	response chan *actor
}

type registerActorWithResponse struct {
	name     string
	response chan *actor
}
