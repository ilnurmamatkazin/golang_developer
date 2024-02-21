package log

import (
	"log"
	"os"
)

// Переопределяем стандартный лог
var (
	Info  = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	Error = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Llongfile)
)
