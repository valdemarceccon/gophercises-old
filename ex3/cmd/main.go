package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/valdemarceccon/gophercises/ex3/ownstory"
)

func main() {
	port := flag.Int("port", 3000, "The port to start CYOS server")
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

	handler := ownstory.NewHandler(story)
	fmt.Printf("Starting server at: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), handler))
}
