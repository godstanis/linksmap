package parser

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// GetContent() func returns actual page content
func TestHttpPage_GetContent(t *testing.T) {
	expectedContent := "<html><h>Our expected resulting content!</h></html>"

	// Mock http server to return our result
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) { rw.Write([]byte(expectedContent)) }))
	defer server.Close()

	Client = server.Client()
	resultContent, err := HttpPage{Path: server.URL}.GetContent()

	if expectedContent != resultContent {
		t.Logf("Expected response:'%s', actual response:'%s'. Error: '%s'", expectedContent, resultContent, err)
		t.Fail()
	}
}

// GetContent() correctly checks http client errors
func TestHttpPage_GetContent_Error(t *testing.T) {
	resultContent, err := HttpPage{Path: "incorrect path"}.GetContent()

	if resultContent != "" || err == nil {
		t.Logf("Incorrect path error expected")
		t.Fail()
	}
}

// GetPath() returns struct's Path
func TestHttpPage_GetPath(t *testing.T) {
	httpPage := HttpPage{Path: "http://example.com"}

	if httpPage.GetPath() != "http://example.com" {
		t.Logf("Incorrect path returned by GetPath() method")
		t.Fail()
	}
}

// GetBasePath() result test
func TestHttpPage_GetBasePath(t *testing.T) {
	testData := []map[string]string{
		{"input": "some_error_path", "output": "some_error_path"},
		{"input": "http://example.com", "output": "http://example.com"},
		{"input": "https://example.com", "output": "https://example.com"},
		{"input": "http://example.com/", "output": "http://example.com"},
		{"input": "http://example.com/something/special", "output": "http://example.com"},
	}

	for _, expected := range testData {
		httpPage := HttpPage{Path: expected["input"]}
		if httpPage.GetBasePath() != expected["output"] {
			t.Logf("Incorrect path returned by GetBasePath('%s') => '%s'", expected["input"], httpPage.GetBasePath())
			t.Fail()
		}
	}
}
