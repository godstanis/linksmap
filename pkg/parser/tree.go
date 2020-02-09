package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type Link struct {
	Value    string `json:"value"`
	Depth    int    `json:"tree_level"`
	Width    int    `json:"tree_width"`
	Children []Link `json:"children"`
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
		panic(e)
	}
}

// Generates tree of links for url with boundaries
func ConstructTreeForUrl(url string, maxWidth int, maxDepth int) Link {
	var baseNode = Link{url, 0, 0, nil}

	fmt.Println("Generating links map for '" + url + "', it can take some time...")

	var wg sync.WaitGroup
	ConstructLinksTreeForNode(&baseNode, maxWidth, maxDepth, 0, &wg)
	wg.Wait()
	return baseNode
}

// Writes json representation of a tree to a file
func DropTreeToJsonFile(tree Link) {
	b, err := json.MarshalIndent(tree, "", "  ")
	check(err)

	dir, err := os.Getwd()
	check(err)

	err = ioutil.WriteFile(dir+"/linksMap.json", b, 0644)
	check(err)
}

// Count all elements in tree recursively
func CountElements(node Link) int {
	var count = 1 // 1 is the firs one
	for _, element := range node.Children {
		if element.Children == nil {
			// Leaf
			count++
			continue
		}
		count += CountElements(element)
	}
	return count
}

// Parse and cunstruct a tree map of urls from main node
func ConstructLinksTreeForNode(node *Link, limitWidth int, limitDepth int, curDepth int, wg *sync.WaitGroup) {
	fmt.Print(".") // Simple output indicator of progress
	links, _ := GetLinksWithAdapter(AdapterResolver{}.GetAdapter(node.Value))
	curDepth++
	if curDepth > limitDepth {
		return
	}

	// It's important to define our cap of node children to prevent 're-referencing' later
	if len(links) <= limitWidth+1 {
		node.Children = make([]Link, len(links))
	} else {
		links = links[:limitWidth] // If links count is more than our limit - cut extra off
		node.Children = make([]Link, limitWidth)
	}

	// Append our links to the node
	for idx, link := range links {
		wg.Add(1)
		node.Children[idx] = Link{link, curDepth, idx, nil}
		go func(idx int) {
			ConstructLinksTreeForNode(&node.Children[idx], limitWidth, limitDepth, curDepth, wg)
			wg.Done()
		}(idx)
	}
}

// Retrieve all the links via a SchemaAdapter
func GetLinksWithAdapter(adapter SchemaAdapter) ([]string, error) {
	var newLinks []string
	body, err := adapter.Content()
	if err != nil {
		log.Fatal(err)
		return newLinks, err
	}

	return SimpleParser{String: body}.ParseLinks(adapter.GetBasePath(), false)
}
