package logging

import "log"

// Logger type
type Logger struct {
	name string
}

// New generate a new logger
func New(s string) *Logger {
	l := Logger{name: s}
	return &l
}

// Message logs a message with using the make of the logger
func (l Logger) Message(s string) {
	log.Printf("[%s] %s\n", l.name, s)
}
