package parser

import (
	"errors"
	"testing"
)

const MockContent = `
		<html>
			<p>Hello world!</p>
			<a href='http://example.com'></a>
			<a href='http://example.com/docs'></a>
			<a href='google.com/news'></a>
			href="not link"
			<a href="http://google.com/news"></a>
			<a href="http://example.com/something"></a>
			href='http://nonono.com'
		</html>
`

// ParseLinks() returns expected parsed links
func TestSimpleParser_ParseLinks(t *testing.T) {
	// "content" : "expectedResult"
	testData := map[string][]string{
		MockContent: {"http://example.com", "http://example.com/docs", "http://example.com/something", "google.com/news", "http://google.com/news"},
		`<html>
			<p>Hello world!</p>
			href="not link"
			href='http://nonono.com'
		</html>`: {},
	}

	for content, expectedLinks := range testData {
		links, _ := simpleParser{content}.ParseLinks()
		if !same(links, expectedLinks) {
			t.Logf("\nExpected links: %s\nActual links:%s", expectedLinks, links)
			t.Fail()
		}
	}
}

// ParsedLinksWithFilters() filters and transforms results using provided functions
func TestSimpleParser_ParseLinksWithFilters(t *testing.T) {
	// "content" : "expectedResult"
	testData := map[string][]string{
		MockContent: {"http://example.com_transformed", "http://example.com/docs_transformed", "http://example.com/something_transformed", "google.com/news_transformed"},
		`<html>
			<p>Hello world!</p>
			href="not link"
			href='http://nonono.com'
		</html>`: {},
	}

	filterFunc := func(link string) bool {
		return link == "http://example.com" || link == "google.com/news"
	}

	transformFunc := func(link string) string {
		return link + "_transformed"
	}

	for content, expectedLinks := range testData {
		links, _ := simpleParser{content}.ParseLinksWithFilters(filterFunc, transformFunc, true)
		if !same(links, expectedLinks) {
			t.Logf("\nExpected links: %s\nActual links:%s", expectedLinks, links)
			t.Fail()
		}
	}
}

// ParseLinksWithFilters() sameDomain parameter returns links only with different domains
func TestSimpleParser_ParseLinksWithFilters_Domain(t *testing.T) {
	// "content" : "expectedResult"
	testData := map[string][]string{
		MockContent: {"http://example.com", "google.com/news", "http://google.com/news"}, // "google.com/news"(relative) and "http://google.com/news" are valid
		`<html>
			<p>Hello world!</p>
			href="not link"
			href='http://nonono.com'
		</html>`: {},
	}

	nullFilter := func(link string) bool { return true }
	nullTransformer := func(link string) string { return link }

	for content, expectedLinks := range testData {
		links, _ := simpleParser{content}.ParseLinksWithFilters(nullFilter, nullTransformer, false)
		if !same(links, expectedLinks) {
			t.Logf("\nExpected links: %s\nActual links:%s", expectedLinks, links)
			t.Fail()
		}
	}
}

// baseUrl simple helper error result test
func TestSimpleParser_baseUrl_err(t *testing.T) {
	if baseUrl("1234") != "1234" {
		t.Logf("baseUrl method is expected to return original link on error")
		t.FailNow()
	}
}

type MockPage struct {
	FakePath     string
	FakeBasePath string
	FakeContent  string
	ContentError error
}

// GetPath Returns absolute path of a page
func (Adapter MockPage) GetPath() string {
	return Adapter.FakePath
}

// GetBasePath Returns base path of a page
func (Adapter MockPage) GetBasePath() string {
	return Adapter.FakeBasePath
}

// GetContent Returns page content
func (Adapter MockPage) GetContent() (string, error) {
	return Adapter.FakeContent, Adapter.ContentError
}

func TestSimpleParser_getLinksForPage_Error(t *testing.T) {
	expectedError := errors.New("mocked error")
	value, err := getLinksForPage(MockPage{
		FakePath:     "http://example.com/test",
		FakeBasePath: "http://example.com",
		FakeContent:  MockContent,
		ContentError: expectedError,
	})

	if !same(value, []string{}) || err != expectedError {
		t.Logf("Mocked error has not been returned from getLinksForPage method!")
		t.Fail()
	}
}

func TestSimpleParser_getLinksForPage(t *testing.T) {
	value, err := getLinksForPage(MockPage{FakePath: "http://example.com/test", FakeBasePath: "http://example.com", FakeContent: MockContent})
	expectedLinks := []string{"http://example.com", "http://google.com/news"}

	if !same(value, expectedLinks) || err != nil {
		t.Log(value, err)
		t.Fail()
	}
}

// Little helper to compare two string slices without index
func same(x, y []string) bool {
	xMap := make(map[string]int)
	yMap := make(map[string]int)

	for _, xElem := range x {
		xMap[xElem]++
	}
	for _, yElem := range y {
		yMap[yElem]++
	}

	for xMapKey, xMapVal := range xMap {
		if yMap[xMapKey] != xMapVal {
			return false
		}
	}
	return true
}
