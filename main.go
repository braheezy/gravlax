package main

import (
	"fmt"
	"os"

	"github.com/braheezy/gravlax/internal/lox"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: gravlax [filename]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		lox.RunFile(os.Args[1])
	} else {
		lox.RunPrompt()
	}
}
