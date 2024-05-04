package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var (
	hadError bool
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: gravlax [filename]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runFile(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	scanner := Scanner{
		source: string(file),
		line:   1,
	}

	run(&scanner)
	if hadError {
		os.Exit(65)
	}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	scanner := Scanner{line: 1}
	for {
		if !scanner.inBlockComment {
			fmt.Print("> ")
		}
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("\nbye!")
				break
			}
			log.Fatal(err)
		}
		scanner.source = strings.TrimSpace(line) // Update source for the new line
		scanner.current = 0                      // Reset current position for new input
		scanner.tokens = nil                     // Clear previous tokens

		run(&scanner)
		hadError = false
	}
}

func run(scanner *Scanner) {

	scanner.scanTokens()
	if !scanner.inBlockComment {
		for _, token := range scanner.tokens {
			fmt.Println("token:", token.String())
		}
	}
}

func reportError(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s:%s\n", line, where, message)
	hadError = true
}
