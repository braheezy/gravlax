package lox

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func RunFile(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	scanner := Scanner{
		Source: string(file),
		Line:   1,
	}

	err = run(&scanner)
	if err != nil {
		if _, ok := err.(*RuntimeError); ok {
			os.Exit(70)
		}
		os.Exit(65)
	}
}

func RunPrompt() {
	reader := bufio.NewReader(os.Stdin)
	scanner := Scanner{Line: 1}
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

		if err := run(&scanner); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	}
}

func run(scanner *Scanner) error {
	err := scanner.ScanTokens()
	if err != nil {
		return err
	}

	if !scanner.InBlockComment {
		parser := Parser{Tokens: scanner.Tokens}
		expression, err := parser.Parse()
		if err != nil {
			return err
		}
		if expression != nil {
			interpret(expression)
		}
	}
	return nil
}

func reportError(lint int, message string) {
	report(lint, "", message)
}

func reportTokenError(token Token, message string) {
	if token.Type == EOF {
		report(token.Line, " at end", message)
	} else {
		report(token.Line, " at '"+token.Lexeme+"'", message)
	}
}

func report(lint int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", lint, where, message)
}
