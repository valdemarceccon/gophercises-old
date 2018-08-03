package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/valdemarceccon/gophercises/ex4/linkparser"
)

func main() {
	htmlFile := flag.String("file", "ex1.html", "File to parse")
	flag.Parse()
	arquivo, err := os.Open(*htmlFile)

	if err != nil {
		fmt.Printf("Could not load the file %s. Error: %s\n", *htmlFile, err)
		os.Exit(1)
	}

	links, err := linkparser.Parse(arquivo)

	if err != nil {
		fmt.Printf("Could not parse the file %s. Error: %s\n", *htmlFile, err)
		os.Exit(1)
	}

	fmt.Println(links)
}
