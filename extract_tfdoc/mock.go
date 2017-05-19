package main

import (
	"fmt"
	"github.com/drhodes/golorem"
	"math/rand"
	"path/filepath"
	"strings"

	data "github.com/coveo/terraform-auto-snippets/common_data"
)

func getArgs() []data.Argument {
	result := make([]data.Argument, rand.Intn(10)+1)
	for i := 0; i < len(result); i++ {
		result[i] = data.Argument{
			Name:        lorem.Word(3, 15),
			Description: lorem.Sentence(2, 10),
			Required:    rand.Intn(3) != 0,
		}
	}
	return result
}

func getResources(path string) []data.Resource {
	result := make([]data.Resource, rand.Intn(200)+3)
	for i := 0; i < len(result); i++ {
		name := lorem.Word(3, 15)
		result[i] = data.Resource{
			Name:        name,
			Description: lorem.Sentence(2, 10),
			URL:         getURL(filepath.Join(path, name)),
			Arguments:   getArgs(),
		}
	}
	return result
}

func getData(path string) []data.Data {
	result := make([]data.Data, rand.Intn(10))
	for i := 0; i < len(result); i++ {
		name := lorem.Word(3, 15)
		result[i] = data.Data{
			Name:        name,
			Description: lorem.Sentence(2, 10),
			URL:         getURL(filepath.Join(path, name)),
			Arguments:   getArgs(),
		}
	}
	return result
}

func getURL(path string) string {
	s := lorem.Sentence(2, 4)
	return fmt.Sprintf("https://www.terraform.io/docs/providers/%s/%s", path, strings.Replace(s[:len(s)-1], " ", "/", -1))
}
