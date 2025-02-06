package logger

import "log"

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

func Error(message string, v ...any) {
	if ERROR >= level {
		log.Printf("[ERROR] %s %v", message, v)
	}
}

func Errorf(message string, v ...any) {
	if ERROR >= level {
		log.Printf("[ERROR] "+message, v...)
	}
}

func Info(message string, v ...any) {
	if INFO >= level {
		log.Printf("[INFO] %s %v", message, v)
	}
}

func Infof(message string, v ...any) {
	if INFO >= level {
		log.Printf("[INFO] "+message, v...)
	}
}

func Debug(message string, v ...any) {
	if DEBUG >= level {
		log.Printf("[DEBUG] %s %v", message, v)
	}
}
func Debugf(message string, v ...any) {
	if DEBUG >= level {
		log.Printf("[DEBUG] "+message, v...)
	}
}

func Warn(message string, v ...any) {
	if WARN >= level {
		log.Printf("[WARN] %s %v", message, v)
	}
}

func Warnf(message string, v ...any) {
	if WARN >= level {
		log.Printf("[WARN] "+message, v...)
	}
}

func Fatal(message string, v ...any) {
	log.Printf("[FATAL] %s %v", message, v)
}

func Fatalf(message string, v ...any) {
	log.Printf("[FATAL] "+message, v...)
}
