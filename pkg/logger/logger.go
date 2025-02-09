package logger

import (
	"log"
)

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	OFF
)

var level int

func Init(l int) {
	level = l
}

func messageToSlice(prefix, message string, slice []any) []any {
	v := make([]any, 2, len(slice)+2)
	v[0] = prefix
	v[1] = message
	v = append(v, slice...)
	return v
}

func Error(message string, v ...any) {
	if ERROR >= level {
		log.Println(messageToSlice("[ERROR]", message, v)...)
	}
}

func Errorf(message string, v ...any) {
	if ERROR >= level {
		log.Printf("[ERROR] "+message, v...)
	}
}

func Info(message string, v ...any) {
	if INFO >= level {
		log.Println(messageToSlice("[INFO]", message, v)...)
	}
}

func Infof(message string, v ...any) {
	if INFO >= level {
		log.Printf("[INFO] "+message, v...)
	}
}

func Debug(message string, v ...any) {
	if DEBUG >= level {
		msg := messageToSlice("[DEBUG]", message, v)
		log.Println(msg...)
	}
}
func Debugf(message string, v ...any) {
	if DEBUG >= level {
		log.Printf("[DEBUG] "+message, v...)
	}
}

func Warn(message string, v ...any) {
	if WARN >= level {
		log.Println(messageToSlice("[WARN]", message, v)...)
	}
}

func Warnf(message string, v ...any) {
	if WARN >= level {
		log.Printf("[WARN] "+message, v...)
	}
}

func Fatal(message string, v ...any) {
	log.Fatal(messageToSlice("[FATAL]", message, v)...)
}

func Fatalf(message string, v ...any) {
	log.Fatalf("[FATAL] "+message, v...)
}
