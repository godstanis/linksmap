package parser

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type Link struct {
	Value    string `json:"value"`
	Id       int    `json:"id"`
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
func ConstructTreeForUrl(url string, maxWidth int, maxDepth int) (Link, error) {
	var baseNode = Link{url, 0, 0, 0, nil}

	log.Println("Generating links map for '" + url + "', it can take some time...")

	var wg sync.WaitGroup
	var step int // Will use reference for cheap 'loose' indexing of tree elements
	err := ConstructLinksTreeForNode(&baseNode, maxWidth, maxDepth, 0, &wg, &step)
	wg.Wait()

	log.Printf("\n Elements in tree: %d\n", CountElements(baseNode))

	log.Println("Links map for '" + url + "' is ready...")
	return baseNode, err
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
func ConstructLinksTreeForNode(node *Link, limitWidth int, limitDepth int, curDepth int, wg *sync.WaitGroup, step *int) error {
	curDepth++
	if curDepth > limitDepth {
		return nil
	}

	links, err := GetLinksWithAdapter(AdapterResolver{}.GetAdapter(node.Value))
	if err != nil || len(links) == 0 {
		return err
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
		log.Println(*step)
		*step++

		node.Children[idx] = Link{link, *step, curDepth, idx, nil}

		go func(idx int) {
			ConstructLinksTreeForNode(&node.Children[idx], limitWidth, limitDepth, curDepth, wg, step)
			wg.Done()
		}(idx)
	}
	return nil
}

// Retrieve all the links via a SchemaAdapter
func GetLinksWithAdapter(adapter SchemaAdapter) ([]string, error) {
	var newLinks []string
	body, err := adapter.Content()
	if err != nil {
		return newLinks, err
	}

	return SimpleParser{String: body}.ParseLinks(adapter.GetBasePath(), false)
}
