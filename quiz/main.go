package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

var quizCorrect, quizNum int

func main() {
	// parse command line argument
	var fileName string
	var timeout string
	var shuffle bool
	flag.StringVar(&fileName, "f", "problem.csv", "CSV file path")
	flag.StringVar(&timeout, "t", "30s", "Duration of the quiz")
	flag.BoolVar(&shuffle, "s", false, "Whether shuffle the quiz list")
	flag.Parse()

	// read CSV file
	var quizList = readQuiz(fileName)
	quizNum = len(quizList)
	if shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(quizNum, func(i, j int) {
			quizList[i], quizList[j] = quizList[j], quizList[i]
		})
	}

	// quiz begin
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Quiz duration: %v. Press ENTER to continue", timeout)
	reader.ReadString('\n')

	// set timeout
	duration, _ := time.ParseDuration(timeout)
	go time.AfterFunc(duration, timeoutFunc)

	// continuosly feed quiz to user
	var ans string
	for _, quiz := range quizList {
		ans = ""
		fmt.Printf("%v? ", quiz[0])
		ans, _ = reader.ReadString('\n')
		ans = strings.Trim(ans, " \r\n")

		fmt.Printf("Your answer: %v\n", ans)
		if ans == quiz[1] {
			quizCorrect += 1
		}
	}
	fmt.Printf("Result: %v/%v\n", quizCorrect, quizNum)
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

func timeoutFunc() {
	fmt.Println("\nTIMEOUT")
	fmt.Printf("Result: %v/%v\n", quizCorrect, quizNum)
	os.Exit(0)
}
