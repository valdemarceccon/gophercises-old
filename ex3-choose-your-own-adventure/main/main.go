package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type history struct {
	arc     string
	content struct {
		title   string
		story   []string
		options []struct {
			text string
			arc  string
		}
	}
}

const jsonFilePath = "/home/valdemar/gocode/src/github.com/valdemarceccon/gophercises/ex3-choose-your-own-adventure/gopher.json"

func main() {

	jsonContent, err := ioutil.ReadFile(jsonFilePath)

	if err != nil {
		fmt.Printf("Cannot open json file: %s. Error: %s", jsonFilePath, err)
		os.Exit(1)
	}

	var hist history

	if err := json.Unmarshal(jsonContent, &hist); err != nil {
		fmt.Printf("Cannot decode json file: %s. Error: %s", jsonFilePath, err)
		os.Exit(1)
	}

	fmt.Println(hist)
}
