package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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

type Information struct {
	Name        string
	Description string
}

type Provider struct {
	Information
	Arguments     []Argument
	DataResources []Data
	Resources     []Resource
}

type Argument struct {
	Information
	Required bool
}

type Data struct {
	Information
}

type Resource struct {
	Information
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

	doc.Find(".active ul a").Each(func(i int, s *goquery.Selection) {
		name := s.Text()
		href, ok := s.Attr("href")
		if ok {
			provider := Provider{
				Information: Information{
					Name:        strings.ToLower(name),
					Description: fmt.Sprintf(`This is the description of "%s"`, name),
				},
				Arguments:     []Argument{},
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
	err = saveToYaml("test.yml", result)
	return
}

func saveToYaml(filename string, data interface{}) error {
	buffer, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, buffer, 0644)
}
