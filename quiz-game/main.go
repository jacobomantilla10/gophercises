package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	filenameFlag := flag.String("filename", "problems.csv", "Name of problems file")
	timerFlag := flag.Int("time", 10, "Length of time that you want the quiz to run for")
	shuffleFlag := flag.Bool("shuffle", false, "Option to shuffle quiz questions around")

	flag.Parse()

	file, err := os.Open(*filenameFlag)

	if err != nil {
		panic(err)
	}

	r := csv.NewReader(file)
	records, err := r.ReadAll()

	if err != nil {
		panic(err)
	}

	if *shuffleFlag {
		records = shuffle(records)
	}

	var enter string
	fmt.Println("Press Enter to start quiz...")
	fmt.Scanln(enter)
	fmt.Println("-----------------------------------------------------------------------------------")
	fmt.Println("Starting Quiz...")
	fmt.Println("-----------------------------------------------------------------------------------")

	timer := time.NewTimer(time.Second * time.Duration(*timerFlag))
	numRight := 0

	for _, record := range records {
		question, answer := strings.TrimSpace(record[0]), strings.TrimSpace(strings.ToLower(record[1]))

		fmt.Fprintf(os.Stdout, "%s: ", question)

		answerCh := make(chan string)

		go func() {
			var input string
			fmt.Scanln(&input)

			answerCh <-strings.ToLower(input)
		}()

		select {
		case <-timer.C:
			fmt.Println("\n-----------------------------------------------------------------------------------")
			fmt.Fprintf(os.Stdout, "Quiz over. Right Answers: %d. Wrong Answers: %d. Total Score: %d/%d\n", 
					numRight, len(records)-numRight, numRight, len(records))
			fmt.Println("-----------------------------------------------------------------------------------")
			return
		case ans := <-answerCh:
			if ans == answer {
				numRight++
			}
		}
	}

	fmt.Println("-----------------------------------------------------------------------------------")
	fmt.Fprintf(os.Stdout, "Quiz over. Right Answers: %d. Wrong Answers: %d. Total Score: %d/%d\n", 
					numRight, len(records)-numRight, numRight, len(records))
	fmt.Println("-----------------------------------------------------------------------------------")
}

func shuffle(records [][]string) [][]string {
	length := len(records)

	for i := range records {
		j := rand.Intn(length)
		records[i], records[j] = records[j], records[i]
	}

	return records
}