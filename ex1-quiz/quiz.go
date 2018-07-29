package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

func main() {
	var file = flag.String("file", "problems.csv", "CSV file containing the quiz problems.")

	flag.Parse()
	problems := parseFile(*file)

	var score uint

	for i, problem := range problems {
		var resp string
		fmt.Printf("%d - %s: ", i+1, problem.question)
		fmt.Scanf("%s\n", &resp)

		if resp == problem.response {
			score++
		}
	}

	fmt.Printf("Score: %d\n", score)

}

func parseFile(filePath string) []problem {
	file, err := os.Open(filePath)

	if err != nil {
		fmt.Printf("Cannot open file: %s\n", filePath)
		os.Exit(1)
	}

	var csvReader = csv.NewReader(file)
	csvContent, err := csvReader.ReadAll()

	result := make([]problem, len(csvContent))

	if err != nil {
		fmt.Printf("Cannot parse csv file: %s. Cause: %s\n", filePath, err)
		os.Exit(1)
	}

	for i, csvline := range csvContent {
		result[i] = problem{
			question: string(csvline[0]),
			response: string(csvline[1]),
		}
	}

	return result
}

type problem struct {
	question string
	response string
}
