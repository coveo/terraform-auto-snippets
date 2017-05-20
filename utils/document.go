package utils

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
)

// GetDocument returns a goquery.Document that could be used to parse
// the content of an HTML file
func GetDocument(uri url.URL) (result *goquery.Document, err error) {
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
