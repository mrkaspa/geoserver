package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mrkaspa/geoserver/handler"
	"github.com/mrkaspa/geoserver/models"
	"github.com/mrkaspa/geoserver/utils"
)

func main() {
	initMain()
	startServer()
}

func startServer() {
	http.Handle("/", handler.NewRouter())
	port := os.Getenv("PORT")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func initMain() {
	utils.LoadEnv(".env_dev")
	utils.InitLogger()
	models.InitDB()
	handler.InitSearcher(models.NewPersistor)
}
