package config

import (
	"log"
	"os"
)

// In Go, an identifier that starts with a capital letter is exported from the package,
// and can be accessed by anyone outside the package that declares it it.
// If an identifier starts with a lower case letter, it can only be accessed from within the package
type Application struct {
	DebugLog *log.Logger
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Port     string // New field for port
	API_key  string // New field for port
}

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
