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

			// if link.String() != "/docs/providers/aws/index.html" {
			// 	return
			// }
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

	id, ok := title.Attr("id")
	if !ok {
		fmt.Fprintf(os.Stderr, "No id found in title for %s\n", uri.String())
		id = strings.ToLower(strings.Replace(titleText, " ", "-", -1))
	}

	provider = &data.Provider{
		Name:          strings.Replace(id, "-provider", "", 1),
		Title:         titleText,
		Description:   title.Next().Text(),
		URL:           uri.String(),
		Arguments:     getArgs(titleText, uri, doc.Find("#argument-reference")),
		DataResources: getResources(titleText, uri, doc.Find("ul.nav.docs-sidenav li.active"), true),
		Resources:     getResources(titleText, uri, doc.Find("ul.nav.docs-sidenav li.active"), false),
	}

	return
}

func getDocument(uri url.URL) (result *goquery.Document, err error) {
	// fmt.Printf("Fetching %s\n", uri.String())
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

func getArgs(name string, uri url.URL, head *goquery.Selection) (arguments []data.Argument) {
	// fmt.Println("Get arguments for", name)
	if head.Length() == 0 {
		fmt.Fprintln(os.Stderr, "No argument for", name)
		return
	}

	var subArgument string
	_ = subArgument
	head.NextFilteredUntil("ul, p", "h2").Each(func(i int, section *goquery.Selection) {
		nodeType := section.Nodes[0].Data
		switch nodeType {
		case "ul":
			section.Find("li").Each(func(i int, arg *goquery.Selection) {
				const req = "(Required)"
				const opt = "(Optional)"
				argName, _ := arg.Find("a").First().Attr("name")
				description := strings.TrimSpace(strings.TrimPrefix(arg.Text(), argName+" - "))
				required := strings.Contains(description, req)
				if required {
					description = strings.Replace(strings.TrimSpace(description), req, "", 1)
				} else {
					description = strings.Replace(strings.TrimSpace(description), opt, "", 1)
				}

				url, _ := arg.Find("a[href]").Attr("href")
				arguments = append(arguments, data.Argument{
					Name:        argName,
					Description: trim(description),
					Required:    required,
					URL:         uri.String() + url,
				})
			})
		case "p":
		}
	})
	return
}

func getResources(name string, parent url.URL, head *goquery.Selection, data bool) (resources []data.Resource) {
	head.SiblingsFiltered("li").Each(func(i int, section *goquery.Selection) {
		sectionTitle := section.ChildrenFiltered("a").Text()
		elements := section.Find("ul li")
		if elements.Length() == 0 {
			return
		}

		const dataSources = "Data Sources"
		var sectionType string
		if data {
			if sectionTitle != dataSources {
				return
			}
			sectionType = "Data"
		} else {
			if sectionTitle == dataSources {
				return
			}
			sectionType = "Resource"
		}

		fmt.Println(sectionType, "-▶", sectionTitle, elements.Length())
		elements.Find("a").Each(func(i int, element *goquery.Selection) {
			href, _ := element.Attr("href")
			link, err := url.Parse(href)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Malformed URL %v for provider %s\n", href, name)
				return
			}
			resource, err := getResource(*parent.ResolveReference(link))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while getting %s: %v\n", href, err)
				return
			}
			resources = append(resources, *resource)
		})
	})

	return
}

func getResource(uri url.URL) (resource *data.Resource, err error) {
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

	id, ok := title.Attr("id")
	if !ok {
		fmt.Fprintf(os.Stderr, "No id found in title for %s\n", uri.String())
		id = strings.ToLower(strings.Replace(titleText, " ", "-", -1))
	}

	resource = &data.Resource{
		Name:        id,
		Description: trim(title.Next().Text()),
		URL:         uri.String(),
		Arguments:   getArgs(titleText, uri, doc.Find("#argument-reference")),
	}

	return
}

func trim(s string) string {
	s = strings.Replace(s, "\n", " ", -1)
	return strings.TrimSpace(s)
}
