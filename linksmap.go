package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/godstanis/linksmap/pkg/parser"
)

func main() {
	http.HandleFunc("/", handler)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("ui/assets/"))))

	log.Println("Server is listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handle incoming http requests
func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "application/json")

		_ = r.ParseMultipartForm(1000)
		pageURL, width, depth := parseFormParameters(r.Form)
		tree, _ := parser.ConstructTreeForUrl(pageURL, width, depth)

		var linksMapJSON, _ = json.MarshalIndent(tree, "", "    ")
		_, _ = w.Write(linksMapJSON)
		return
	}
	http.ServeFile(w, r, "./ui/index.html") // serve static index.html
}

// Parse required form parameters
func parseFormParameters(form url.Values) (string, int, int) {
	pageURL := ""
	maxWidth, _ := strconv.ParseInt(form.Get("width"), 10, 64)
	maxDepth, _ := strconv.ParseInt(form.Get("depth"), 10, 64)

	if maxWidth == 0 || maxWidth > 6 {
		maxWidth = 2
	}

	if maxDepth == 0 || maxDepth > 6 {
		maxDepth = 2
	}

	if pageURL = form.Get("url"); pageURL != "" && !strings.HasPrefix(pageURL, "http") {
		pageURL = "http://" + pageURL
	}

	return pageURL, int(maxWidth), int(maxDepth)
}
