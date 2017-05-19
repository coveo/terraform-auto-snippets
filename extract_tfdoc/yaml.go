package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func saveToYaml(filename string, data interface{}) error {
	buffer, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, buffer, 0644)
}
