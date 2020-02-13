package parser

import (
	"log"
	"sync"
)

// Link represents a parsed site URL, it is a node,
// children are links found on parent node's html page response.
type Link struct {
	Value    string   `json:"value"`    // Full url
	Info     LinkInfo `json:"info"`     // Additional information object
	Children []Link   `json:"children"` // Slice of Links found on this link's url
}

// LinkInfo represents some additional information regarding link (it's position, index and so on)
type LinkInfo struct {
	Id         int    `json:"id"`          // Unique id for the node in a tree
	ShortValue string `json:"value_short"` // Shorthand value/name
	Depth      int    `json:"depth"`       // Depth of the node in a tree
	Width      int    `json:"width"`       // Position of node in children slice (0 for each first node child)
}

// Generates tree of links for url with boundaries
func ConstructTreeForUrl(url string, maxWidth int, maxDepth int) (Link, error) {
	log.Println("Generating links map for '" + url + "'...")

	var baseNode = Link{url, LinkInfo{ShortValue: url}, nil}

	var wg sync.WaitGroup
	err := ConstructLinksTreeForNode(&baseNode, maxWidth, maxDepth, 0, &wg)
	wg.Wait()
	// Populating our tree with additional data before returning
	populateTreeInfo(&baseNode, 0)

	log.Printf("\nTree for %s is generated.\n Elements count: %d\n", url, countElementsInTree(baseNode))
	return baseNode, err
}

// Parses and constructs a tree map of urls from main node
func ConstructLinksTreeForNode(node *Link, limitWidth int, limitDepth int, curDepth int, wg *sync.WaitGroup) error {
	curDepth++
	if curDepth >= limitDepth {
		return nil
	}

	links, err := getLinksForPage(HttpPage{node.Value})
	if err != nil || len(links) == 0 {
		return err
	}

	// It's important to define our cap of node children to prevent 're-referencing' later
	if len(links) < limitWidth {
		node.Children = make([]Link, len(links))
	} else {
		links = links[:limitWidth] // If links count is more than our limit - cut extra off
		node.Children = make([]Link, limitWidth)
	}

	// Append our links to the node
	for idx, link := range links {
		wg.Add(1)

		node.Children[idx] = Link{
			link,
			LinkInfo{Depth: curDepth, Width: idx},
			nil,
		}

		go func(idx int) {
			_ = ConstructLinksTreeForNode(&node.Children[idx], limitWidth, limitDepth, curDepth, wg)
			wg.Done()
		}(idx)
	}
	return nil
}

// Set unique indexes to each node and update additional information
func populateTreeInfo(node *Link, index int) int {
	node.Info.Id = index
	node.Info.ShortValue = baseUrl(node.Value)
	index++
	for idx := range node.Children {
		index = populateTreeInfo(&node.Children[idx], index)
	}
	return index
}

// Count all elements in tree recursively
func countElementsInTree(node Link) int {
	var count = 1 // 1 is the firs one
	for _, element := range node.Children {
		if element.Children == nil {
			count++
			continue
		}
		count += countElementsInTree(element)
	}
	return count
}
