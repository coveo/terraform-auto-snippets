package utils

import (
	"fmt"
	"os"
)

// TrapErrors recovers any error raised by panic and print it to stderr
func TrapErrors(handler func(string, ...interface{})) {
	if err := recover(); err != nil {
		if handler == nil {
			handler = func(format string, args ...interface{}) {
				PrintError("%v", err)
				os.Exit(1)
			}
		}
		handler("%v", err)
	}
}

// Assert issues a panic if the condition is not met
func Assert(condition bool, format string, args ...interface{}) {
	if !condition {
		panic(fmt.Errorf(format, args...))
	}
}

// PanicOnError issues a panic if there is an error
func PanicOnError(err error, args ...interface{}) {
	if err != nil {
		var userMessage string
		if len(args) > 0 {
			userMessage = fmt.Sprintf(" "+fmt.Sprintf("%v", args[0]), args[1:]...)
		}
		panic(fmt.Errorf("Error %v%s", err, userMessage))
	}
}
