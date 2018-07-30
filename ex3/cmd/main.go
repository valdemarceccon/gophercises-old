package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/valdemarceccon/gophercises/ex3/ownstory"
)

func main() {
	filePath := flag.String("file", "gopher.json", "File path to story JSON.")
	file, err := os.Open(*filePath)

	if err != nil {
		fmt.Printf("Cannot open json file: %s. Error: %s", *filePath, err)
		os.Exit(1)
	}

	story, err := ownstory.LoadStoryJSON(file)

	if err != nil {
		fmt.Printf("Cannot decode json file: %s. Error: %s", *filePath, err)
		os.Exit(1)
	}

	fmt.Println(story)
}
