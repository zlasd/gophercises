package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

func main() {
	// read csv file and gen quiz
	var fileName string
	flag.StringVar(&fileName, "f", "problem.csv", "CSV file path")
	flag.Parse()
	var quizList = readQuiz(fileName)

	// continuosly feed quiz to user
	var correct int = 0
	var ans string = ""
	for _, quiz := range quizList {
		fmt.Printf("%v? ", quiz[0])
		fmt.Scanf("%v\n", &ans)
		fmt.Printf("Your answer: %v\n", ans)
		if ans == quiz[1] {
			correct += 1
		}
	}
	fmt.Printf("%v/%v\n", correct, len(quizList))
}

func readQuiz(fileName string) [][]string {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Open file %v failed: %v", fileName, err)
		os.Exit(-1)
	}
	bufReader := bufio.NewReader(file)
	csvReader := csv.NewReader(bufReader)
	records, err := csvReader.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "file %v cannot parsed as a CSV file", fileName)
		os.Exit(-1)
	}
	return records
}
