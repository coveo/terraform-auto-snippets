package main

import (
	"fmt"
	"os"
)

func main() {
	providers, err := getTerraformDocs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	err = saveToYaml("mock.yml", providers)
}
