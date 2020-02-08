package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type Link struct {
	Value    string `json:"value"`
	Depth    int    `json:"tree_level"`
	Width    int    `json:"tree_width"`
	Children []Link `json:"children"`
}

type AdapterResolver struct{}

func (Adapter AdapterResolver) GetAdapter(path string) SchemaAdapter {
	if strings.HasPrefix(path, "http") {
		return HttpSchemaAdapter{path}
	}

	return FileSchemaAdapter{path}
}

type SchemaAdapter interface {
	Content() (string, error)
	GetPath() string
	GetBasePath() string
}

type FileSchemaAdapter struct {
	Path string
}

func (Adapter FileSchemaAdapter) GetPath() string {
	return Adapter.Path
}
func (Adapter FileSchemaAdapter) GetBasePath() string {
	split := strings.Split(Adapter.GetPath(), "/")
	split = split[:len(split)-1]
	return strings.Join(split, "/") + "/"
}
func (Adapter FileSchemaAdapter) Content() (string, error) {
	body, err := ioutil.ReadFile(Adapter.GetPath())
	if err != nil {
		log.Fatal(err)
	}
	return string(body), err
}

type HttpSchemaAdapter struct {
	Path string
}

func (Adapter HttpSchemaAdapter) GetPath() string {
	return Adapter.Path
}
func (Adapter HttpSchemaAdapter) GetBasePath() string {
	return Adapter.GetPath()
}
func (Adapter HttpSchemaAdapter) Content() (string, error) {
	client := http.Client{Timeout: 15 * time.Second}
	response, err := client.Get(Adapter.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	return string(body), err
}

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

func check(e error) {
	if e != nil {
		log.Fatal(e)
		panic(e)
	}
}

func main() {
	url := "https://www.google.ru/"
	// url := "/home/garstas/Development/html/urlmap_test/index.html"
	var width, depth = 3, 4
	var parent = Link{url, 0, 0, nil}

	fmt.Println("Generating links map for '" + url + "', it can take some time...")
	ConstructLinksTreeForNode(&parent, width, depth, 0)

	DropTreeToJsonFile(parent)
	fmt.Printf("\nElements: %d", CountElements(parent))
	fmt.Println("\nFile has been generated. Have fun! :-D")
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
func ConstructLinksTreeForNode(node *Link, limitWidth int, limitDepth int, curDepth int) {
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
		node.Children[idx] = Link{link, curDepth, idx, nil}
		ConstructLinksTreeForNode(&node.Children[idx], limitWidth, limitDepth, curDepth)
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

	return SimpleParser{body}.ParseLinks(adapter.GetBasePath(), false)
}
