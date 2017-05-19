package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"os"
	"strings"

	data "github.com/coveo/terraform-auto-snippets/common_data"
)

const TerraformDocURL = "https://www.terraform.io/docs/providers/index.html"

func getTerraformHelp() (err error) {
	uri, err := url.Parse(TerraformDocURL)
	if err != nil {
		return err
	}

	providers, err := getProviders(*uri)
	if err != nil {
		return err
	}

	return saveToYaml("../mock.yml", providers)
}

func getProviders(uri url.URL) (providers map[string]data.Provider, err error) {
	providers = map[string]data.Provider{}
	doc, err := getDocument(uri)
	if err != nil {
		return
	}

	doc.Find(".active ul a").Each(func(i int, s *goquery.Selection) {
		name := s.Text()
		href, ok := s.Attr("href")
		if ok {
			link, err := url.Parse(href)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Malformed URL %v for provider %s\n", href, name)
				return
			}

			if provider, err := getProvider(*uri.ResolveReference(link)); err == nil {
				providers[name] = *provider
			} else {
				fmt.Fprintf(os.Stderr, "Unable to get provider %s: %v\n", name, err)
			}
		}
	})

	return
}

func getProvider(uri url.URL) (provider *data.Provider, err error) {
	doc, err := getDocument(uri)
	if err != nil {
		return
	}

	title := doc.Find("h1")
	switch title.Length() {
	case 0:
		err = fmt.Errorf("No title found in %s", uri.String())
		return
	case 1:
		break
	default:
		fmt.Fprintf(os.Stderr, "Found more that one title (%d) in %s\n", title.Length(), uri.String())
		title = title.First()
	}

	titleNode := title.Nodes[0].LastChild
	if titleNode.Type != html.TextNode {
		err = fmt.Errorf("Malformed title node in %s %v", uri.String(), title.Text())
		return
	}

	titleText := strings.TrimSpace(titleNode.Data)

	// doc.Find("h1").Each(func(i int, s *goquery.Selection) {
	id, ok := title.Attr("id")
	if !ok {
		fmt.Fprintf(os.Stderr, "No id found in title for %s\n", uri.String())
		id = strings.ToLower(strings.Replace(titleText, " ", "-", -1))
	}

	provider = &data.Provider{
		Name:          id,
		Title:         titleText,
		Description:   title.Next().Text(),
		URL:           uri.String(),
		Arguments:     getArgs(),
		DataResources: getData(id),
		Resources:     getResources(id),
	}

	return
}

func getDocument(uri url.URL) (result *goquery.Document, err error) {
	fmt.Printf("Fetching %s\n", uri.String())
	response, err := http.Get(uri.String())
	switch {
	case response == nil:
		return
	case response.StatusCode >= 400:
		err = fmt.Errorf("%s -▶ %s", TerraformDocURL, response.Status)
	}
	if err == nil {
		result, err = goquery.NewDocumentFromResponse(response)
	}
	return
}
