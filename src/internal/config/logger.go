package config

import (
	"log"
	"os"
)

var AppLog *Application

func init() {
	initLogger()
}

func initLogger() {
	AppLog = &Application{
		DebugLog: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime),
		InfoLog:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime),
	}
}
