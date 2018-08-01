package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func main() {
	htmlFile := flag.String("file", "ex1.html", "File to parse")
	flag.Parse()
	arquivo, err := os.Open(*htmlFile)

	if err != nil {
		fmt.Printf("Could not load the file %s. Error: %s\n", *htmlFile, err)
		os.Exit(1)
	}

	z := html.NewTokenizer(arquivo)

	// links := make([]Link, 10)

	for {
		if z.Next() == html.ErrorToken {
			fmt.Printf("%v\n", z.Token())
			break
		}

		v, _ := z.TagName()

		fmt.Printf("%s - %v\n", v, z.Token())
	}

}

type Link struct {
	Href string
	Text string
}
