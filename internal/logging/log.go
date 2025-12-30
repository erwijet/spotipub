package logging

import (
	"log"
	"os"
)

func GetLogger(ctx string) *log.Logger {
	return log.New(os.Stdout, "["+ctx+"] ", log.Ldate|log.Ltime)
}
