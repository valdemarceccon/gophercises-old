package linkparser

import (
	"io"
	"strings"

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

	links := parse(node.FirstChild)

	return links, nil
}

func parse(node *html.Node) []Link {
	var ret []Link
	links := linkNodes(node)

	for _, l := range links {
		ret = append(ret, Link{text(l), href(l)})
	}

	return ret
}

func linkNodes(node *html.Node) []*html.Node {
	var ret []*html.Node

	if node.Data == "a" {
		ret = append(ret, node)
	} else {
		for inner := node.FirstChild; inner != nil; inner = inner.NextSibling {
			ret = append(ret, linkNodes(inner)...)
		}
	}
	return ret
}

func text(node *html.Node) string {
	var ret string
	for inner := node.FirstChild; inner != nil; inner = inner.NextSibling {
		if inner.Type == html.TextNode {
			ret += inner.Data
		} else {
			ret += text(inner)
		}
	}

	ret = strings.Join(strings.Fields(ret), " ")

	return ret
}

func href(node *html.Node) string {
	for _, v := range node.Attr {
		if v.Key == "href" {
			return v.Val
		}
	}

	return ""
}
