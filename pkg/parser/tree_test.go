package parser

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

// Fake html to test on
const MockContentForTree = `
		<html>
			<p>Hello world!</p>
			<a href='http://example.com'></a>
			<a href='http://example.com/docs'></a>
			<a href='example2.com/news'></a>
			href="not link"
			<a href="http://example2.com/news"></a>
			<a href="http://example3.com/something"></a>
			<a href="https://example4.com/something"></a>
			<a href="https://example5.com/something"></a>
			href='http://nonono.com'
		</html>
`

// Fake client Transport
type MockTripper struct {
	ReturnString string
}

func (tt MockTripper) RoundTrip(rq *http.Request) (*http.Response, error) {
	if tt.ReturnString == "" {
		tt.ReturnString = MockContentForTree
	}
	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString(tt.ReturnString)),
	}, nil
}

func TestConstructTreeForUrl(t *testing.T) {
	const maxWidth, maxDepth, expectedCount = 2, 3, 7

	Client = &http.Client{Transport: MockTripper{}}
	tree, err := ConstructTreeForUrl("http://test.test", maxWidth, maxDepth)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if count := countElementsInTree(tree); count != expectedCount {
		t.Logf("Unexpected tree result. Expected %d, got %d", expectedCount, count)
		t.FailNow()
	}

	if tree.Value != "http://test.test" || tree.Info.ShortValue != "test.test" || tree.Info.Width != 0 || tree.Info.Depth != 0 {
		t.Log("Base node is incorrect")
		t.FailNow()
	}

	if len(tree.Children) != maxWidth {
		t.Log("Width is incorrect")
		t.FailNow()
	}

	if len(tree.Children) != maxWidth || len(tree.Children[0].Children) != maxWidth || len(tree.Children[0].Children[0].Children) != 0 {
		t.Log("Depth is incorrect")
		t.FailNow()
	}
}
