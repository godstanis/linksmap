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
)

type Link struct {
	Value    string
	Depth    int
	Width    int
	Children []Link
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
	response, err := http.Get(Adapter.Path)
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
	links := re.FindAllStringSubmatch(string(Parser.String), -1)
	for _, element := range links {
		// if it's relative link - prepend our main url to it

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
	//url := "/home/garstas/Development/html/urlmap_test/index.html"
	var width, depth = 3, 4
	var parent = Link{url, 0, 0, nil}

	fmt.Println("Generating links map for '" + url + "', it can take some time...")
	ConstructLinksTreeForNode(&parent, width, depth, 0)

	DropTreeToJsonFile(parent)
	fmt.Println("File has been generated. Have fun! :-D")
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

func countLinksInNode(node Link) int {
	var count = 1 // 1 is the firs one
	for _, element := range node.Children {
		if element.Children == nil {
			// Leaf
			count++
		} else {
			count += countLinksInNode(element)
		}
	}
	return count
}

func ConstructLinksTreeForNode(node *Link, limitWidth int, limitDepth int, curDepth int) {
	var links, _ = GetLinksWithAdapter(AdapterResolver{}.GetAdapter(node.Value))

	curDepth++
	if curDepth > limitDepth {
		return
	}
	for idx, link := range links {
		if limitWidth <= idx {
			break
		}
		// Add new link child to node
		node.Children = append(node.Children, Link{link, curDepth, idx, nil})

		if curDepth > limitDepth {
			fmt.Println("bleat")
		}
		ConstructLinksTreeForNode(&node.Children[idx], limitWidth, limitDepth, curDepth)
	}
}

// Retrieves all the links via a SchemaAdapter
func GetLinksWithAdapter(adapter SchemaAdapter) ([]string, error) {
	var newLinks []string
	body, err := adapter.Content()
	if err != nil {
		log.Fatal(err)
		return newLinks, err
	}

	parser := SimpleParser{body}

	return parser.ParseLinks(adapter.GetBasePath(), false)
}
