package main

import (
	"fmt"
	"linksmap/pkg/parser"
)

type Link struct {
	Value    string `json:"value"`
	Depth    int    `json:"tree_level"`
	Width    int    `json:"tree_width"`
	Children []Link `json:"children"`
}

func main() {
	var url = "https://www.google.ru/"
	// var url = "/home/garstas/Development/html/urlmap_test/index.html"

	var tree = parser.ConstructTreeForUrl(url, 4, 4)

	parser.DropTreeToJsonFile(tree)
	fmt.Printf("\nElements: %d", parser.CountElements(tree))
	fmt.Println("\nFile has been generated. Have fun! :-D")
}
