package main

import (
	"encoding/json"
	"linksmap/pkg/parser"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc("/", handler)
	log.Println("Server is listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handle incoming http requests
func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "application/json")

		_ = r.ParseMultipartForm(1000)
		pageUrl, width, depth := parseFormParameters(r.Form)
		tree, _ := parser.ConstructTreeForUrl(pageUrl, width, depth)

		var linksMapJson, _ = json.MarshalIndent(tree, "", "    ")
		_, _ = w.Write(linksMapJson)
		return
	}
	http.ServeFile(w, r, "./ui/index.html") // serve static index.html
}

// Parse required form parameters
func parseFormParameters(form url.Values) (string, int, int) {
	pageUrl := ""
	maxWidth, _ := strconv.ParseInt(form.Get("width"), 10, 64)
	maxDepth, _ := strconv.ParseInt(form.Get("depth"), 10, 64)

	if maxWidth == 0 || maxWidth > 6 {
		maxWidth = 2
	}

	if maxDepth == 0 || maxDepth > 6 {
		maxDepth = 2
	}

	if pageUrl = form.Get("url"); pageUrl != "" && !strings.HasPrefix(pageUrl, "http") {
		pageUrl = "http://" + pageUrl
	}

	return pageUrl, int(maxWidth), int(maxDepth)
}
