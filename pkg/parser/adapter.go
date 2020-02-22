package parser

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	Client *http.Client
)

func init() {
	Client = &http.Client{Timeout: 15 * time.Second}
}

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
	parsedUrl, err := url.Parse(Adapter.GetPath())
	if err != nil || !parsedUrl.IsAbs() {
		return Adapter.GetPath()
	}
	return parsedUrl.Scheme + "://" + parsedUrl.Hostname()
}

// GetContent Returns page content
func (Adapter HttpPage) GetContent() (string, error) {
	response, err := Client.Get(Adapter.Path)
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
