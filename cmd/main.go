package main

import (
	"log"
	"net/http"
	"os"
)

var (
	hostname string = ""
)

func init() {
	hostname, _ = os.LookupEnv("LINKSMAP_HOSTNAME")
}

func main() {
	http.HandleFunc(hostname+"/", getBaseHandleFunc())
	http.Handle(hostname+"/assets/", getStaticHandler())

	log.Println("Server is listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
