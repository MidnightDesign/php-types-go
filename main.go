package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

const typesDirPath = "../php-types-rust/tests/types/"

func main() {
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
			tokens := lex(line)
			node := parseNode(newCursor(tokens))
			if line != node.print() {
				fmt.Println("Before: " + line)
				fmt.Println("After:  " + node.print())
			}
		}
		return nil
	})
	tokens := lex("array{name: string, age: number}")
	fmt.Printf("%v\n", tokens)
}
