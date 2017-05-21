package utils

import (
	"fmt"
	"github.com/gruntwork-io/terragrunt/errors"
	"os"
)

// TrapErrors recovers any error raised by panic and print it to stderr
func TrapErrors(handler func(string, ...interface{})) {
	if err := recover(); err != nil {
		var text string
		switch err := err.(type) {
		case error:
			text = errors.PrintErrorWithStackTrace(err)
		default:
			text = fmt.Sprintf("%v", err)
		}
		if handler == nil {
			handler = func(format string, args ...interface{}) {
				PrintError("%v", text)
				os.Exit(1)
			}
		}
		handler("%v", text)
	}
}

// Trap any panic error and add stack trace to the resulting error
func TrapPanic() {
	if err := recover(); err != nil {
		switch err := err.(type) {
		case error:
			panic(errors.WithStackTrace(err))
		default:
			fmt.Println(2)
			panic(errors.WithStackTrace(fmt.Errorf("%v", err)))
		}
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
