package logging

import (
	"log"
)

// Logger type
type Logger struct {
	name    string
	enabled bool
}

// New generate a new logger
func New(s string) *Logger {
	l := Logger{name: s, enabled: false}
	return &l
}

// Message logs a message with using the make of the logger
func (l *Logger) Message(s string, v ...string) {

	if l.enabled {
		if len(v) > 0 {
			for _, value := range v {
				s = s + " " + value
			}
		}

		log.Printf("[%s] %s\n", l.name, s)
	}
}

// Enable ...
func (l *Logger) Enable() {
	l.enabled = true
}

// Disable ...
func (l *Logger) Disable() {
	l.enabled = false
}
