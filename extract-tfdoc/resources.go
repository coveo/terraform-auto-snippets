package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"net/url"
	"strings"
	"sync"

	data "github.com/coveo/terraform-auto-snippets/common_data"
	"github.com/coveo/terraform-auto-snippets/utils"
)

type resourceWorkForce struct {
	ProviderWorkForce
	providerName string
	providerURL  url.URL
	head         *goquery.Selection
}

func (wf resourceWorkForce) String() string {
	return fmt.Sprintf("%s (%s)", wf.providerName, wf.providerURL.String())
}

func (wf resourceWorkForce) get(data bool) (resources []data.Resource) {
	sections := wf.head.SiblingsFiltered("li")
	var selector = func(sel *goquery.Selection) *goquery.Selection { return sel.Find("ul li") }

	var warning = func(format string, args ...interface{}) {
		if data {
			// We only print the warning on Data to avoid printing it twice
			utils.PrintWarning(format, args...)
		}
	}

	switch sections.Length() {
	case 0:
		utils.PrintError("Unable to find resources section for %s", wf)
	case 1:
		// This is a special case where the data is in the current section (occurs for External and HTTP)
		warning("Special handling to find resources for %s", wf)
		sections = selector(wf.head)
		selector = func(sel *goquery.Selection) *goquery.Selection { return sel }
	}

	sections.Each(func(i int, section *goquery.Selection) {
		sectionTitle := section.ChildrenFiltered("a").Text()
		elements := selector(section)

		switch sectionTitle {
		case "All Providers", "« Documentation Home":
			return
		case "":
			warning("Special handling to find section title for %s", wf)
			sectionTitle = section.Parent().Find("h4").Text()
		}

		if elements.Length() == 0 {
			previous := section.PrevFiltered("h4")
			if previous.Length() == 1 {
				sectionTitle = previous.Text()
				if data {
					utils.PrintWarning("Title stored in previous element for %s in %s", sectionTitle, wf)
				}
				elements = section
			} else {
				warning("No element in section %s for %s", sectionTitle, wf)
				return
			}
		}

		if sectionTitle == "" {
			utils.PrintError("Unable to find section name for %s", wf)
			return
		}

		dataSources := strings.Contains(sectionTitle, "Data Source")
		if data && !dataSources || !data && dataSources {
			return
		}

		var (
			mutex    sync.Mutex
			subTasks sync.WaitGroup
		)

		elements.Find("a").Each(func(i int, element *goquery.Selection) {
			href, _ := element.Attr("href")

			subTasks.Add(1)
			go func() {
				defer subTasks.Done()

				resource, err := wf.getResource(sectionTitle, wf.providerURL, href)
				if err != nil {
					utils.PrintError("   Error while getting %s: %v", href, err)
					return
				}

				defer mutex.Unlock()
				mutex.Lock()
				resources = append(resources, *resource)
			}()
		})
		subTasks.Wait()
	})
	return
}

// getResource returns detail about the requested resource
func (wf resourceWorkForce) getResource(section string, parent url.URL, link string) (resource *data.Resource, err error) {
	uri := wf.resolve(&wf.providerURL, link)
	if uri == nil {
		err = fmt.Errorf("Malformed URI %s", link)
		return
	}

	err = wf.ProcessDocument(*uri, func(doc *goquery.Document, responseChannel chan error) {
		var err error
		defer func() { responseChannel <- err }()

		title := doc.Find("h1")
		switch title.Length() {
		case 0:
			err = fmt.Errorf("No title found in %s", uri.String())
			return
		case 1:
			break
		default:
			utils.PrintWarning("Found more that one title (%d) in %s", title.Length(), uri.String())
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
			utils.PrintWarning("No id found in title for %s", uri.String())
			id = strings.ToLower(strings.Replace(titleText, " ", "-", -1))
		}

		resource = &data.Resource{
			Base: data.Base{
				Name:        id,
				Description: utils.Trim(title.Next().Text()),
				URL:         uri.String(),
			},
			Section:   section,
			Arguments: getArgs(titleText, *uri, doc.Find("#argument-reference")),
		}
		return
	})
	return
}
