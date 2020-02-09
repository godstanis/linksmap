package main

import (
	"encoding/json"
	"linksmap/pkg/parser"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)
	log.Println("Server is listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "application/json")

		r.ParseMultipartForm(1000)
		var tree parser.Link
		if r.Form.Get("url") != "" {
			tree, _ = parser.ConstructTreeForUrl(r.Form.Get("url"), 4, 4)
		}

		var json, _ = json.MarshalIndent(tree, "", "    ")
		w.Write(json)
		return
	}
	http.ServeFile(w, r, "./ui/index.html") // serve static index.html
}
