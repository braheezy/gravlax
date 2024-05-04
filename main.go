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

	run(string(file))
	if hadError {
		os.Exit(65)
	}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("\nbye!")
				break
			}
			log.Fatal(err)
		}
		run(strings.TrimSpace(line))
		hadError = false
	}
}

func run(source string) error {
	scanner := Scanner{
		source: source,
		line:   1,
	}
	scanner.scanTokens()

	for _, token := range scanner.tokens {
		fmt.Println("token:", token.String())
	}
	return nil
}

func reportError(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s:%s\n", line, where, message)
	hadError = true
}
