package utils

import (
	"os"

	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/rifflock/lfshook"
)

//LogWriter for logging
var Log *log.Logger

//InitLog system
func InitLogger() {
	if Log == nil {
		log.SetOutput(os.Stderr)
		log.SetLevel(log.WarnLevel)
		Log = log.New()
		Log.Formatter = new(log.JSONFormatter)
		fmt.Printf("LOADED LOG %s", os.Getenv("MODE"))
		if os.Getenv("MODE") == "production" {
			Log.Hooks.Add(lfshook.NewHook(lfshook.PathMap{
				log.InfoLevel:  "log/info.log",
				log.ErrorLevel: "log/error.log",
			}))
		}
	}
}
