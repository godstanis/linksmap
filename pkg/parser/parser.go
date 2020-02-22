// Package parser parses websites and generates tree structures (i.e. maps) of their connections.
package parser

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// Simple parser to parse links from a website.
type simpleParser struct {
	Content string
}

// Parses links without any additional filters
func (Parser simpleParser) ParseLinks() ([]string, error) {
	emptyFilter := func(url string) bool {
		return true
	}
	emptyTransform := func(url string) string {
		return url
	}

	return Parser.ParseLinksWithFilters(emptyFilter, emptyTransform, true)
}

// Parses links with user provided filter and transformer closure.
//
// filter: accepts currently found link as a parameter and returns boolean.
// If the value returned by filter is false than this link will not be included in final result.
//
// transform: accepts currently found link as a parameter and returns transformed result.
//
// allowSameDomain: determines if same domain urls are included
func (Parser simpleParser) ParseLinksWithFilters(filter func(url string) bool, transform func(url string) string, allowSameDomain bool) ([]string, error) {
	parsedLinks := make(map[string][]string)

	doc, err := html.Parse(strings.NewReader(Parser.Content))
	if err != nil {
		return []string{}, err
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					if filter(a.Val) {
						transformedLink := transform(a.Val)

						hostname := baseUrl(transformedLink)

						if _, ok := parsedLinks[hostname]; !ok {
							parsedLinks[hostname] = []string{}
						}
						parsedLinks[hostname] = append(parsedLinks[hostname], transformedLink)
					}
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	var resultLinks []string
	// Transform map to slice
	for _, links := range parsedLinks {
		if !allowSameDomain {
			resultLinks = append(resultLinks, links[0])
		} else {
			resultLinks = append(resultLinks, links...)
		}
	}
	return resultLinks, nil
}

// Return base url for passed link
func baseUrl(link string) string {
	parsedUrl, err := url.Parse(link)
	if err != nil || parsedUrl.Host == "" {
		return link
	}
	return parsedUrl.Hostname()
}

// Retrieve all the links via a Page
func getLinksForPage(page Page) ([]string, error) {
	body, err := page.GetContent()
	if err != nil {
		return []string{}, err
	}

	// Filter all related links off
	urlFilter := func(link string) bool {
		parsedUrl, err := url.Parse(link)
		if err != nil {
			fmt.Println(err)
			return false
		}
		return parsedUrl.IsAbs()
	}

	return simpleParser{body}.ParseLinksWithFilters(urlFilter, func(link string) string { return link }, false)
}
