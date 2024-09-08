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
		"Assign   : name Token, value Expr",
		"Binary   : left Expr, operator Token, right Expr",
		"Call     : callee Expr, paren Token, arguments []Expr",
		"Grouping : expression Expr",
		"Literal  : value interface{}",
		"Logical  : left Expr, operator Token, right Expr",
		"Unary    : operator Token, right Expr",
		"Variable : name Token",
	}, "Eval", "(interface{}, *RuntimeError)")
	if err != nil {
		log.Fatal(err)
	}
	err = defineAst(outputDir, "Stmt", []string{
		"Block      : statements []Stmt",
		"Expression : expression Expr",
		"Function   : name Token, params []Token, body []Stmt",
		"If         : condition Expr, thenBranch Stmt, elseBranch Stmt",
		"Print      : expression Expr",
		"Return     : keyword Token, value Expr",
		"Var        : initializer Expr, name Token",
		"While      : condition Expr, body Stmt",
		"Break      : ",
	}, "Execute", "*RuntimeError")
	if err != nil {
		log.Fatal(err)
	}
}

func defineAst(dir string, baseName string, types []string, methodName string, methodArgs string) error {
	path := path.Join(dir, fmt.Sprintf("%v.go", strings.ToLower(baseName)))
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	writer.WriteString("// auto-generated by generate_ast.go. DO NOT EDIT")
	writer.WriteRune('\n')

	writer.WriteString("package lox\n")
	writer.WriteRune('\n')

	writer.WriteString(fmt.Sprintf("type %v interface {\n", baseName))
	writer.WriteString(fmt.Sprintf("%v() %v\n}", methodName, methodArgs))
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
