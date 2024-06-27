package utils

import (
	"log"
)

func LogInfo(format string, v ...any) {
	log.SetPrefix(ActionStyle.Render("[INFO] "))
	log.Printf(format, v...)
}

func LogWarning(format string, v ...any) {
	log.SetPrefix(WarningStyle.Render("[WARN] "))
	log.Printf(format, v...)
}

func LogError(format string, v ...any) {
	log.SetPrefix(ErrorStyle.Render("[ERRO] "))
	log.Fatalf(format, v...)
}
