package parser

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

// Page represents any parsable object (http page, file e.t.c.)
//
// Any page should provide ability to determine it's content by path
type Page interface {
	GetContent() (string, error)
	GetPath() string
	GetBasePath() string
}

// HttpPage represents web page with html content
type HttpPage struct {
	Path string
}

// GetPath Returns absolute path of a page
func (Adapter HttpPage) GetPath() string {
	return Adapter.Path
}

// GetBasePath Returns base path of a page
func (Adapter HttpPage) GetBasePath() string {
	parsedUrl, _ := url.Parse(Adapter.GetPath())
	return parsedUrl.Host + "://" + parsedUrl.Hostname()
}

// GetContent Returns page content
func (Adapter HttpPage) GetContent() (string, error) {
	client := http.Client{Timeout: 15 * time.Second}
	response, err := client.Get(Adapter.Path)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Println(err)
	}

	return string(body), err
}
