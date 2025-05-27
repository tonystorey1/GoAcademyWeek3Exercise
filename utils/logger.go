package utils

import (
	"log"
	"os"
)

const logFileName = "TodoList.log"

var (
	Logger *log.Logger
)

func SetupLogger() {
	file, err := os.OpenFile(logFileName /*os.O_APPEND|*/, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	Logger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}
