package use

import (
	"bytes"
	"errors"
	"golang.org/x/net/html"
	"io"
	"strings"
)

func NumberOfChildren(n *html.Node) int {
	if n == nil {
		return -1
	}

	count := 0
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		count += 1
	}

	return count
}

func ParseHtml(content string) *html.Node {
	output, err := html.Parse(strings.NewReader(content))
	if err != nil {

		return nil
	}

	return output
}

func Body(doc *html.Node) (*html.Node, error) {
	var body *html.Node
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "body" {
			body = node
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	if body != nil {
		return body, nil
	}
	return nil, errors.New("missing <body> in the node tree")
}

func RenderNode(n *html.Node) (string, error) {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	err := html.Render(w, n)
	return buf.String(), err
}
