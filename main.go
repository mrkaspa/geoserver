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
	host := os.Getenv("HOST")
	err := http.ListenAndServe(host+":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func initMain() {
	utils.LoadEnv(".env")
	utils.InitLogger()
	models.InitDB()
	handler.InitSearcher(models.NewPersistor)
}
