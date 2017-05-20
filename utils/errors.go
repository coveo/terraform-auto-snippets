package utils

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

var (
	infoPrinter    = color.New(color.FgGreen).SprintfFunc()
	warningPrinter = color.New(color.FgYellow).SprintfFunc()
	errorPrinter   = color.New(color.FgRed).SprintfFunc()
)

// PrintInfo is used to print a colored message to the stderr
func PrintInfo(format string, args ...interface{}) {
	fmt.Fprintln(os.Stderr, infoPrinter(format, args...))
}

// PrintWarning is used to print a yellow warning message to the stderr
func PrintWarning(format string, args ...interface{}) {
	fmt.Fprintln(os.Stderr, warningPrinter(format, args...))
}

// PrintError is used to print a red error message to the stderr
func PrintError(format string, args ...interface{}) {
	fmt.Fprintln(os.Stderr, errorPrinter(format, args...))
}
