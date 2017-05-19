package main

import (
	"flag"
<<<<<<< HEAD
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

=======
	"log"
	"os"
	"io/ioutil"
	"fmt"
	"path/filepath"
	"strings"
>>>>>>> cf5507e15182ae023dec12872f9aaa715359c56c
	"gopkg.in/yaml.v2"
)

var (
	importDir = flag.String("import-dir", "", "Directerory with yaml files to import")
<<<<<<< HEAD
	vscodeArg = flag.Bool("vscode", false, "Do the converstion for VsCode")
	atomArg   = flag.Bool("atom", false, "Do the converstion for Atom")
)

func listYaml(dir string) ([]os.FileInfo, error) {

	files, err := ioutil.ReadDir(dir)
=======
	vscode    = flag.Bool("vscode", false, "Do the converstion for VsCode")
	atom      = flag.Bool("atom", false, "Do the converstion for Atom")
)

type Provider struct {
	Name 		string 		`yaml:"name"`
	Description 	string 		`yaml:"description"`
	Arguments	[]Argument 	`yaml:"arguments"`
	Dataresources	[]string	`yaml:"dataresources"`
	Resources	[]string	`yaml:"resources"`
}

type Argument struct {
	Name 		string `yaml:"name"`
	Description 	string `yaml:"description"`
	Requierd	bool 	`yame:"required"`
}

func listYaml(dir string) ([]os.FileInfo, error) {

	files , err := ioutil.ReadDir(dir)
>>>>>>> cf5507e15182ae023dec12872f9aaa715359c56c

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

<<<<<<< HEAD
	if *vscodeArg == false && *atomArg == false {
=======
	if *vscode == false && *atom == false {
>>>>>>> cf5507e15182ae023dec12872f9aaa715359c56c
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
<<<<<<< HEAD
	for _, f := range yamlfiles {
		data, err := ioutil.ReadFile(*importDir + f.Name())

		if err != nil {
=======
	for _, f := range yamlfiles{
		data, err := ioutil.ReadFile(*importDir + f.Name())

		if err !=nil {
>>>>>>> cf5507e15182ae023dec12872f9aaa715359c56c
			log.Fatal(err)
		}

		yaml.Unmarshal(data, &p)
	}
<<<<<<< HEAD
	VscodeCreateSnippets(&p)
=======
	fmt.Println(p)
>>>>>>> cf5507e15182ae023dec12872f9aaa715359c56c
}
