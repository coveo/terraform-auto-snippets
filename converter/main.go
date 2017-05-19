package main

import (
	"flag"

	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	importDir = flag.String("import-dir", "", "Directerory with yaml files to import")
	vscodeArg = flag.Bool("vscode", false, "Do the converstion for VsCode")
	atomArg   = flag.Bool("atom", false, "Do the converstion for Atom")
)

func listYaml(dir string) ([]os.FileInfo, error) {

	files, err := ioutil.ReadDir(dir)

	if err != nil {
		return nil, err
	}

	yamlFile := []os.FileInfo{}
	for _, f := range files {
		ext := strings.ToLower(filepath.Ext(f.Name()))
		if ext == ".yaml" || ext == ".yml" {
			yamlFile = append(yamlFile, f)
		}
	}
	return yamlFile, nil
}

func main() {
	flag.Parse()

	if *vscodeArg == false && *atomArg == false {

		log.Println("You have to specify one of editor for convertion")
		flag.Usage()
		os.Exit(1)
	}

	if *importDir == "" {
		log.Println("You have to specify where are the file")
		flag.Usage()
		os.Exit(1)
	}

	yamlfiles, err := listYaml(*importDir)

	if err != nil {
		log.Fatal(err)
	}
	var p map[string]Provider
	for _, f := range yamlfiles {
		data, err := ioutil.ReadFile(*importDir + f.Name())

		if err != nil {
			log.Fatal(err)
		}

		yaml.Unmarshal(data, &p)
	}
	VscodeCreateSnippets(&p)

}
