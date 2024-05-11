package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/braheezy/gravlax/internal/lox"
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

	scanner := lox.Scanner{
		Source: string(file),
		Line:   1,
	}

	err = run(&scanner)
	if err != nil {
		os.Exit(65)
	}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	scanner := lox.Scanner{Line: 1}
	for {
		if !scanner.InBlockComment {
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
		scanner.Source = strings.TrimSpace(line) // Update source for the new line
		scanner.Current = 0                      // Reset current position for new input
		scanner.Tokens = nil                     // Clear previous tokens

		run(&scanner)
	}
}

func run(scanner *lox.Scanner) error {

	err := scanner.ScanTokens()
	if !scanner.InBlockComment {
		for _, token := range scanner.Tokens {
			fmt.Println("token:", token.String())
		}
	}
	return err
}
