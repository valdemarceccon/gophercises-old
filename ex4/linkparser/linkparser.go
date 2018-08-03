package linkparser

import (
	"fmt"
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

func nodeLinks(node *html.Node) []*html.Node {
	var ret []*html.Node

	if node.Data == "a" {
		ret = append(ret, node)
	} else {
		for nNode := node.FirstChild; nNode != nil; nNode = node.NextSibling {
			if nNode.Data == "a" {
				ret = append(ret, nNode)
			} else {
				ret = append(ret, nodeLinks(nNode)...)
			}
		}
	}

	return ret
}

func teste(node *html.Node) []Link {
	result := []Link{}

	nodelinks := nodeLinks(node)

	for _, nodelink := range nodelinks {
		fmt.Println(nodelink)
	}

	return result
}
