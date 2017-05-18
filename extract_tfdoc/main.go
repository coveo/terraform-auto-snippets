package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/drhodes/golorem"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	// "google.golang.org/api/urlshortener/v1"
	"path/filepath"
)

func main() {
	err := getTerraformHelp()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

type Provider struct {
	Name          string
	Description   string
	URL           string
	Arguments     []Argument
	DataResources []Data
	Resources     []Resource
}

type Argument struct {
	Name        string
	Description string
	Required    bool
}

type Data struct {
	Name        string
	Description string
	URL         string
	Arguments   []Argument
}

type Resource Data

func getTerraformHelp() (err error) {
	uri := "https://www.terraform.io/docs/providers/index.html"
	response, err := http.Get(uri)
	switch {
	case response == nil:
		return
	case response.StatusCode >= 400:
		err = fmt.Errorf("%s -▶ %s", uri, response.Status)
	}
	if err != nil {
		return
	}

	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return
	}

	result := map[string]Provider{}

	doc.Find(".active ul a").Each(func(i int, s *goquery.Selection) {
		name := s.Text()
		href, ok := s.Attr("href")
		if ok {
			provider := Provider{
				Name:          strings.ToLower(name),
				Description:   fmt.Sprintf(`This is the description of "%s"`, name),
				URL:           getURL(name),
				Arguments:     getArgs(),
				DataResources: getData(name),
				Resources:     getResources(name),
			}
			result[name] = provider
			fmt.Println(i, name)
			switch name {
			case "AWS":
				base, _ := url.Parse(uri)
				link, err := url.Parse(href)
				if err != nil {
					return
				}
				fmt.Printf("    %s\n", base.ResolveReference(link))
				//		getProviderInfo(doc)
			}
		}
	})

	err = saveToYaml("../mock.yml", result)
	return
}

func saveToYaml(filename string, data interface{}) error {
	buffer, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, buffer, 0644)
}

func getArgs() []Argument {
	result := make([]Argument, rand.Intn(10)+1)
	for i := 0; i < len(result); i++ {
		result[i] = Argument{
			Name:        lorem.Word(3, 15),
			Description: lorem.Sentence(2, 10),
			Required:    rand.Intn(3) != 0,
		}
	}
	return result
}

func getResources(path string) []Resource {
	result := make([]Resource, rand.Intn(200)+3)
	for i := 0; i < len(result); i++ {
		name := lorem.Word(3, 15)
		result[i] = Resource{
			Name:        name,
			Description: lorem.Sentence(2, 10),
			URL:         getURL(filepath.Join(path, name)),
			Arguments:   getArgs(),
		}
	}
	return result
}

func getData(path string) []Data {
	result := make([]Data, rand.Intn(10))
	for i := 0; i < len(result); i++ {
		name := lorem.Word(3, 15)
		result[i] = Data{
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
