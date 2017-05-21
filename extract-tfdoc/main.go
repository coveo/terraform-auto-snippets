package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
	"os"

	"github.com/coveo/terraform-auto-snippets/utils"
)

var (
	app       = kingpin.New(os.Args[0], "Yaml extrator from terraform documentation web site.")
	filters   = app.Flag("filter", "Filter the providers to those containing the supplied parameter").Short('f').PlaceHolder("provider").Strings()
	outfile   = app.Flag("out-file", "Output the result to the specified file").Short('o').OpenFile(os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	nbWorkers = app.Flag("workers", fmt.Sprintf("Number of worker to start (1 to %d)", maxWorkers)).Short('w').Default("5").Int()
)

const maxWorkers = 100

func main() {
	// Handle eventual panic message
	defer utils.TrapErrors(app.Fatalf)

	app.Author("Coveo")
	kingpin.CommandLine = app
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	utils.Assert(*nbWorkers >= 1, "The number of workers must be greater than 1")
	utils.Assert(*nbWorkers <= maxWorkers, "The maximum number of workers is %d", maxWorkers)

	if *outfile == nil {
		utils.PrintWarning("The result will go to stdout")
		*outfile = os.Stdout
	}

	process()
}

func process() {
	wf := ParseTerraformDocumentation(*nbWorkers)
	defer wf.TerminateAll()

	providers, err := wf.GetProviders(*filters...)
	utils.PanicOnError(err, "while parsing terraform documentation web site")

	buffer, err := yaml.Marshal(providers)
	utils.PanicOnError(err, "while converting the result to YAML")

	_, err = (*outfile).Write(buffer)
	utils.PanicOnError(err, "while writing the output file")
}
