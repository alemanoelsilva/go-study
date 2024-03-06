package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

const EXIT = 1

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit("Failed to open the CSV file")
	}

	csvReader := csv.NewReader(file)
	lines, err := csvReader.ReadAll()
	if err != nil {
		exit("Failed to part the provided CSV file.")
	}

	problems := parseLine(lines)

	totalScore := 0

	for index, problem := range problems {
		fmt.Printf("Problem #%d: %s = ", index+1, problem.question)

		timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
		answerChannel := make(chan string)

		/**
		 * This code defines an anonymous goroutine using the go keyword. Goroutines are lightweight threads of execution in Go.
		 * Inside the goroutine, it reads input from the user using fmt.Scanf() and sends the input to the answerChannel.
		 * answerChannel is a channel of strings that is used to communicate the user's answer to the main goroutine.
		 *
		 * Using a go routine:
		 *  - Non-blocking User Input:
		 *     - Using a goroutine allows the program to continue executing other tasks while waiting for user input.
		 *     - Without a goroutine, the program would block until the user provides input, preventing other tasks
		 *       (such as monitoring the timer) from executing concurrently.
		 *  - Avoiding Timer Blocking:
		 *     - Executing fmt.Scanf() directly within the select statement would block the timer from triggering its case
		 *       until the user provides input.
		 *     - By using a goroutine, the timer can continue running independently, ensuring that the program can respond
		 *       to both user input and timer expiration events concurrently.
		 *  - Preventing Deadlocks:
		 *     - Using blocking calls for user input within the select statement could lead to deadlocks if the user does not
		 *       provide input within the specified time limit.
		 *     - Using a goroutine ensures that the program remains responsive and can handle both user input and timer events gracefully,
		 *       preventing potential deadlock situations.
		 */
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerChannel <- answer
		}()

		// The select statement waits for one of its cases to be ready to proceed.
		select {
		/**
		 * The first case <-timer.C waits for the timer to expire.
		 * When the timer's channel (timer.C) receives a message (i.e., when the timer expires), it triggers this case.
		 *  - This indicates that the user did not provide an answer within the time limit.
		 *    In this case, it calls showScore() to display the current score and exits the program.
		 */
		case <-timer.C:
			showScore(totalScore, len(problems))
			return
		/**
		 * The second case <-answerChannel waits for input from the user through the answerChannel.
		 * When the user submits an answer, it receives the answer from the channel.
		 *  - It compares the received answer with the correct answer (problem.answer).
		 *  - If they match, it increments the totalScore counter, indicating that the user answered correctly.
		 */
		case answer := <-answerChannel:
			if answer == problem.answer {
				totalScore++
			}
		}
	}

	showScore(totalScore, len(problems))
}

type problem struct {
	question string
	answer   string
}

func parseLine(lines [][]string) []problem {
	problems := make([]problem, len(lines))
	for index, line := range lines {
		problems[index] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return problems
}

func showScore(totalScore int, total int) {
	fmt.Printf("\nYou scored %d out of %d\n", totalScore, total)
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(EXIT)
}
