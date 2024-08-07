package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

const typesDirPath = "../php-types-rust/tests/types/"
const invalidFilePath = "../php-types-rust/tests/invalid.txt"

func main() {
	f, err := os.Open(invalidFilePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	typeString := ""
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " ")
		if line != "" {
			typeString += line + "\n"
			continue
		}
		tokens, err := newLexer(typeString).lex()
		if err != nil {
			typeString = ""
			continue
		}
		node := parseNode(newCursor(tokens))
		if node == nil {
			typeString = ""
			continue
		}
		fmt.Println("Expected an error for:")
		fmt.Println(typeString)
		fmt.Println("But valid type:")
		fmt.Println(node.print())
	}

	root := os.DirFS(typesDirPath)
	_ = fs.WalkDir(root, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		f, err := root.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.Trim(scanner.Text(), " ")
			if line == "" {
				continue
			}
			tokens, err := newLexer(line).lex()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			node := parseNode(newCursor(tokens))
			if line != node.print() {
				fmt.Println("Before: " + line)
				fmt.Println("After:  " + node.print())
			}
		}
		return nil
	})
}
