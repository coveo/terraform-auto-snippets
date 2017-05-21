package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"

	data "github.com/coveo/terraform-auto-snippets/common_data"
	"github.com/coveo/terraform-auto-snippets/utils"
)

var (
	app       = kingpin.New(os.Args[0], "Terraform extension snippet generator for VSCode and Atom.")
	prefix    = app.Flag("prefix", "Prefix to put before snippet name").Short('p').Default("tf").String()
	vscodeArg = app.Flag("vscode", "Do the conversion for VsCode").Short('v').Bool()
	atomArg   = app.Flag("atom", "Do the conversion for Atom").Short('a').Bool()
	shortSnip = app.Flag("short", "Indicates to generate long version of snippets (enabled by default, --no-short to disable)").Default("true").Bool()
	longSnip  = app.Flag("long", "Indicates to generate short version of snippets (enabled by default, --no-long to disable)").Default("true").Bool()
	fullSnip  = app.Flag("full", "Indicates to generate full version of snippets").Short('f').Default("false").Bool()
	files     = app.Arg("file", "Yaml files to import").ExistingFiles()
)

func main() {
	// Handle eventual panic message
	defer utils.TrapErrors(app.Fatalf)

	app.Author("Coveo")
	kingpin.CommandLine = app
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	utils.Assert(*vscodeArg || *atomArg, "You must specify at least one editor for convertion")
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		*files = append(*files, "[stdin]")
	}

	utils.Assert(len(*files) > 0, "You must specify at least one file to import")

	process()
}

func process() {
	var providers data.ProviderList

	for _, file := range *files {
		var buffer []byte
		var err error

		switch file {
		case "[stdin]":
			buffer, err = ioutil.ReadAll(os.Stdin)
		default:
			buffer, err = ioutil.ReadFile(file)
		}
		utils.PanicOnError(err)

		var p data.ProviderList
		err = yaml.Unmarshal(buffer, &p)
		utils.PanicOnError(err, "Error while decoding %s", file)
		providers = append(providers, p...)
	}

	if *vscodeArg {
		VscodeCreateSnippets(providers)
	}
	if *atomArg {
		utils.PrintWarning("Atom exporter not implemented yet")
	}
}
