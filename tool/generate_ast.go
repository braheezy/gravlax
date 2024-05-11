package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: generate_ast.go <output directory>")
		os.Exit(64)
	}
	outputDir := os.Args[1]

	err := defineAst(outputDir, "Expr", []string{
		"Binary : left Expr, operator Token, right Expr",
		"Grouping : expression Expr",
		"Literal  : value interface{}",
		"Unary    : operator Token, right Expr",
	})
	if err != nil {
		log.Fatal(err)
	}
}

func defineAst(dir string, baseName string, types []string) error {
	path := path.Join(dir, fmt.Sprintf("%v.go", strings.ToLower(baseName)))
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	writer.WriteString("package lox\n")
	writer.WriteRune('\n')

	writer.WriteString(fmt.Sprintf("type %v interface {}\n", baseName))
	writer.WriteRune('\n')

	for _, astType := range types {
		typeName := strings.TrimSpace(strings.Split(astType, ":")[0])
		fields := strings.TrimSpace(strings.Split(astType, ":")[1])
		defineType(writer, typeName, fields)
	}
	return nil
}

func defineType(writer *bufio.Writer, typeName string, fieldList string) {
	writer.WriteString(fmt.Sprintf("type %v struct {\n", typeName))
	for _, field := range strings.Split(fieldList, ",") {
		field = strings.TrimSpace(field)
		writer.WriteString(fmt.Sprintf("  %v\n", field))
	}
	writer.WriteString("}\n")
	writer.WriteRune('\n')
}
