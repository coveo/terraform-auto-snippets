package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"strings"

	data "github.com/coveo/terraform-auto-snippets/common_data"
	"github.com/fatih/color"
	"os"
)

func getDocument(uri url.URL) (result *goquery.Document, err error) {
	response, err := http.Get(uri.String())
	switch {
	case response == nil:
		return
	case response.StatusCode >= 400:
		err = fmt.Errorf("%s -▶ %s", uri.String(), response.Status)
	}
	if err == nil {
		result, err = goquery.NewDocumentFromResponse(response)
	}
	return
}

func getArgs(name string, uri url.URL, head *goquery.Selection) (arguments []data.Argument) {
	if head.Length() == 0 {
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

func trim(s string) string {
	s = strings.Replace(s, "\n", " ", -1)
	return strings.TrimSpace(s)
}

var (
	infoPrinter    = color.New(color.FgGreen).SprintfFunc()
	warningPrinter = color.New(color.FgYellow).SprintfFunc()
	errorPrinter   = color.New(color.FgRed).SprintfFunc()
)

// PrintInfo is used to print a colored message to the stderr
func PrintInfo(format string, args ...interface{}) {
	fmt.Fprintln(os.Stderr, infoPrinter(format, args...))
}

// PrintWarning is used to print a yellow warning message to the stderr
func PrintWarning(format string, args ...interface{}) {
	fmt.Fprintln(os.Stderr, warningPrinter(format, args...))
}

// PrintError is used to print a red error message to the stderr
func PrintError(format string, args ...interface{}) {
	fmt.Fprintln(os.Stderr, errorPrinter(format, args...))
}
