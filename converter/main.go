package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

var (
	app       = kingpin.New(os.Args[0], "Terraform extension snippet generator for VSCode and Atom.")
	vscodeArg = app.Flag("vscode", "Do the conversion for VsCode").Short('v').Bool()
	atomArg   = app.Flag("atom", "Do the conversion for Atom").Short('a').Bool()
	files     = app.Arg("file", "Yaml files to import").ExistingFiles()
)

func main() {
	if err := process(); err != nil {

	}
}

func process() error {
	app.Author("Coveo")
	kingpin.CommandLine = app
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	if *vscodeArg == false && *atomArg == false {
		return fmt.Errorf("You must specify at least one editor for convertion")
	}

	if len(*files) == 0 {
		return fmt.Errorf("You have to specify at least one file to import")
	}

	var p map[string]Provider
	for _, file := range *files {
		data, err := ioutil.ReadFile(file)

		if err != nil {
			log.Fatal(err)
		}

		yaml.Unmarshal(data, &p)
	}

	VscodeCreateSnippets(&p)
	return nil
}
