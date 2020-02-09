package parser

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

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
		log.Println(err)
		return "", err
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
