package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
	"os"
)

var (
	app       = kingpin.New(os.Args[0], "Yaml extrator from terraform documentation web site.")
	filters   = app.Flag("filter", "Filter the providers to those containing the supplied parameter").Short('f').PlaceHolder("provider").Strings()
	outfile   = app.Flag("out-file", "Output the result to the specified file").Short('o').OpenFile(os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	nbWorkers = app.Flag("workers", "Number of worker to start").Short('w').Default("10").Int()
)

func main() {
	app.Author("Coveo")
	kingpin.CommandLine = app
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	if *outfile == nil {
		PrintWarning("The result will go to stdout")
		*outfile = os.Stdout
	}

	if err := process(); err != nil {
		PrintError("%v", err)
	}
}

func process() (err error) {
	providers, err := getTerraformDocs(*filters...)
	if err != nil {
		return fmt.Errorf("%v while parsing terraform documentation web site", err)
	}

	buffer, err := yaml.Marshal(providers)
	if err != nil {
		return fmt.Errorf("%v while converting the result to YAML", err)
	}

	_, err = (*outfile).Write(buffer)
	if err != nil {
		return fmt.Errorf("%v while writing the output file", err)
	}
	return
}
