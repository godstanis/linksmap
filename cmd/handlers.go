package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/godstanis/linksmap/pkg/parser"
)

// Get index http handler function
func getBaseHandleFunc() func(w http.ResponseWriter, r *http.Request) {
	return handleFunc
}

// Get static assets handler
func getStaticHandler() http.Handler {
	return http.StripPrefix("/assets/", http.FileServer(http.Dir("ui/assets/")))
}

// Handle incoming http requests
func handleFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "application/json")

		r.ParseMultipartForm(1000)
		pageURL, width, depth := parseFormParameters(r.Form)
		tree, _ := parser.ConstructTreeForUrl(pageURL, width, depth)

		var linksMapJSON, _ = json.MarshalIndent(tree, "", "    ")
		w.Write(linksMapJSON)
		return
	}
	http.ServeFile(w, r, "./ui/index.html") // serve static index.html
}

// Parse required form parameters
func parseFormParameters(form url.Values) (string, int, int) {
	pageURL := ""
	maxWidth, err := strconv.ParseInt(form.Get("width"), 10, 64)
	if err != nil || maxWidth == 0 || maxWidth > 6 {
		maxWidth = 2
	}
	maxDepth, err := strconv.ParseInt(form.Get("depth"), 10, 64)
	if err != nil || maxDepth == 0 || maxDepth > 6 {
		maxWidth = 2
	}

	if pageURL = form.Get("url"); pageURL != "" && !strings.HasPrefix(pageURL, "http") {
		pageURL = "http://" + pageURL
	}

	return pageURL, int(maxWidth), int(maxDepth)
}
