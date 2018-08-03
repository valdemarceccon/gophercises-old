package linkparser

import (
	"io"

	"golang.org/x/net/html"
)

// Link structure that will hold result after
// parse document
type Link struct {
	Text string
	Href string
}

// Parse html
func Parse(r io.Reader) ([]Link, error) {
	node, err := html.Parse(r)

	if err != nil {
		return nil, err
	}

	result := teste(node)

	return result, nil
}

func text(node *html.Node) string {

	if node.Type == html.TextNode {
		return node.Data
	}

	var result string

	for nodesib := node.FirstChild; nodesib != nil; nodesib = nodesib.NextSibling {
		result += text(nodesib) + " "
	}

	return result
}

func teste(node *html.Node) []Link {
	result := []Link{}
	if node.Type == html.ElementNode && node.Data == "a" {
		var link Link
		link.Text = text(node)
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				link.Href = attr.Val
			}
		}

		result = append(result, link)
	}
	for c := node.FirstChild; c != nil; c = node.NextSibling {
		result = append(result, teste(c)...)
	}

	return result
}
