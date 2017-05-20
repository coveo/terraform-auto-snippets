package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"

	"github.com/coveo/terraform-auto-snippets/utils"
)

var (
	app       = kingpin.New(os.Args[0], "Terraform extension snippet generator for VSCode and Atom.")
	vscodeArg = app.Flag("vscode", "Do the conversion for VsCode").Short('v').Bool()
	atomArg   = app.Flag("atom", "Do the conversion for Atom").Short('a').Bool()
	files     = app.Arg("file", "Yaml files to import").ExistingFiles()
)

func main() {
	if err := process(); err != nil {
		utils.PrintError("%v", err)
		app.Usage(os.Args[1:])
	}
}

func process() (err error) {
	app.Author("Coveo")
	kingpin.CommandLine = app
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	if *vscodeArg == false && *atomArg == false {
		return fmt.Errorf("You must specify at least one editor for convertion")
	}

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		*files = append(*files, "-")
	}

	if len(*files) == 0 {
		return fmt.Errorf("You have to specify at least one file to import")
	}

	var p map[string]Provider
	for _, file := range *files {
		var data []byte
		if file == "-" {
			data, err = ioutil.ReadAll(os.Stdin)
			file = "stdin"
		} else {
			data, err = ioutil.ReadFile(file)
		}
		if err != nil {
			log.Fatal(err)
		}

		err = yaml.Unmarshal(data, &p)
		if err != nil {
			log.Fatal(fmt.Errorf("Error when decoding %s: %v", file, err))
		}
	}

	VscodeCreateSnippets(p)
	return
}
