package parser

import (
	"regexp"
	"strings"
)

type SimpleParser struct {
	String string
}

func (Parser SimpleParser) ParseLinks(basePath string, skipSameBase bool) ([]string, error) {
	var newLinks []string
	re := regexp.MustCompile(`(?m)<a\s+(?:[^>]*?\s+)?href="([^"]*)"`)
	matches := re.FindAllStringSubmatch(string(Parser.String), -1)
	for _, element := range matches {
		// if it's relative link - prepend our base path to it
		var relative = !strings.HasPrefix(element[1], basePath) && !strings.HasPrefix(element[1], "http")
		if relative {
			element[1] = basePath + element[1]
		}
		if skipSameBase && relative {
			continue
		}
		newLinks = append(newLinks, element[1])
	}
	return newLinks, nil
}
