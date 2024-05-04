package config

import (
	"log"
)

// In Go, an identifier that starts with a capital letter is exported from the package,
// and can be accessed by anyone outside the package that declares it it.
// If an identifier starts with a lower case letter, it can only be accessed from within the package
type Application struct {
	DebugLog *log.Logger
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Port     string
	API_key  string
}
