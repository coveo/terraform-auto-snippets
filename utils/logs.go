package utils

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"sync"
)

var (
	MessagePrinter = color.New(color.Reset).SprintfFunc()
	InfoPrinter    = color.New(color.FgGreen).SprintfFunc()
	WarningPrinter = color.New(color.FgYellow).SprintfFunc()
	ErrorPrinter   = color.New(color.FgRed).SprintfFunc()
)

// PrintMessage is used to print a message to the stderr
func PrintMessage(format string, args ...interface{}) { print(MessagePrinter(format, args...)) }

// PrintInfo is used to print a colored message to the stderr
func PrintInfo(format string, args ...interface{}) { print(InfoPrinter(format, args...)) }

// PrintWarning is used to print a yellow warning message to the stderr
func PrintWarning(format string, args ...interface{}) { print(WarningPrinter(format, args...)) }

// PrintError is used to print a red error message to the stderr
func PrintError(format string, args ...interface{}) { print(ErrorPrinter(format, args...)) }

func print(messageformat string) {
	defer printMutex.Unlock()
	printMutex.Lock()
	fmt.Fprintln(os.Stderr, messageformat)
}

var printMutex sync.Mutex
