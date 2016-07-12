package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/mrkaspa/matchserver/models"
	"github.com/mrkaspa/matchserver/utils"
	"github.com/mrkaspa/matchserver/ws"
)

func main() {
	initMain()
	startServer()
}

func startServer(){
	router := mux.NewRouter()
	router.HandleFunc("/ws/{username}", ws.ServeWS)
	http.Handle("/", router)
	port := os.Getenv("PORT")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func initMain() {
	if err := godotenv.Load(".env_dev"); err != nil {
		log.Fatal("Error loading .env_dev file")
	}
	utils.Init()
	models.Init()
}
