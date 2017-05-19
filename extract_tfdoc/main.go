package main

import (
	"fmt"
	"os"
)

func main() {
	err := getTerraformHelp()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
