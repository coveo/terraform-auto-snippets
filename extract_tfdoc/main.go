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
}

type Resource struct {
	Name        string
	Description string
}

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

	r := rand.New(rand.NewSource(0))

	doc.Find(".active ul a").Each(func(i int, s *goquery.Selection) {
		name := s.Text()
		href, ok := s.Attr("href")
		if ok {
			arguments := make([]Argument, r.Intn(10)+1)
			for i := 0; i < len(arguments); i++ {
				arguments[i] = Argument{
					Name:        lorem.Word(3, 15),
					Description: lorem.Sentence(2, 10),
					Required:    r.Intn(3) != 0,
				}
			}
			provider := Provider{
				Name:          strings.ToLower(name),
				Description:   fmt.Sprintf(`This is the description of "%s"`, name),
				Arguments:     arguments,
				DataResources: []Data{},
				Resources:     []Resource{},
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

	fmt.Println(result)
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
