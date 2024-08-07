package main

import (
	"fmt"
	"unicode"
)

type tokenType int

const (
	identifier tokenType = iota + 1
	integer
	stringLiteral
	openAngle
	closeAngle
	openCurly
	closeCurly
	openBrace
	closeBrace
	comma
	colon
	pipe
	equal
	hyphen
	question
	ellipsis
	ampersand
)

var singleCharTokens = map[rune]tokenType{
	'<': openAngle,
	'>': closeAngle,
	'{': openCurly,
	'}': closeCurly,
	'(': openBrace,
	')': closeBrace,
	',': comma,
	':': colon,
	'|': pipe,
	'=': equal,
	'-': hyphen,
	'?': question,
	'&': ampersand,
}

type token struct {
	t           tokenType
	stringValue string
}

func isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isNumeric(c rune) bool {
	return c >= '0' && c <= '9'
}

func isIdentifierChar(c rune) bool {
	return isAlpha(c) || isNumeric(c) || c == '-'
}

func parseIdentifierToken(c *cursor[rune]) token {
	char, eof := c.peek()
	if eof {
		return token{t: identifier, stringValue: ""}
	}
	idStr := string(char)
	c.next()
	for {
		char, eof = c.peek()
		if eof {
			break
		}
		if !isIdentifierChar(char) {
			break
		}
		idStr += string(char)
		c.next()
	}
	return token{t: identifier, stringValue: idStr}
}

func parseInt(c *cursor[rune]) token {
	char, eof := c.peek()
	if eof {
		return token{t: integer, stringValue: ""}
	}
	numStr := string(char)
	c.next()
	for {
		char, eof = c.peek()
		if eof {
			break
		}
		if !isNumeric(char) {
			break
		}
		numStr += string(char)
		c.next()
	}
	return token{t: integer, stringValue: numStr}
}

func skipWhitespace(c *cursor[rune]) {
	for {
		char, eof := c.peek()
		if eof {
			break
		}
		if !unicode.IsSpace(char) {
			break
		}
		c.next()
	}
}

func parseStringLiteralToken(c *cursor[rune]) token {
	c.next()
	str := ""
	for {
		char, eof := c.peek()
		if eof {
			break
		}
		if char == '\'' {
			c.next()
			break
		}
		str += string(char)
		c.next()
	}
	return token{t: stringLiteral, stringValue: str}
}

func parseEllipsis(c *cursor[rune]) token {
	periods := 0
	for {
		char, eof := c.peek()
		if eof {
			break
		}
		if char != '.' {
			break
		}
		periods++
		c.next()
		if periods != 3 {
			continue
		}
		break
	}
	return token{t: ellipsis, stringValue: "..."}
}

func lex(s string) []token {
	var chars []rune
	for _, char := range s {
		chars = append(chars, char)
	}
	cursor := newCursor(chars)
	tokens := make([]token, 0)
	for {
		char, eof := cursor.peek()
		if eof {
			break
		}
		tokenType := singleCharTokens[char]
		if tokenType != 0 {
			tokens = append(tokens, token{t: tokenType, stringValue: string(char)})
			cursor.next()
			continue
		}
		if char == '.' {
			tokens = append(tokens, parseEllipsis(cursor))
			continue
		}
		if char == '\'' {
			tokens = append(tokens, parseStringLiteralToken(cursor))
			continue
		}
		if unicode.IsSpace(char) {
			skipWhitespace(cursor)
			continue
		}
		if isAlpha(char) {
			tokens = append(tokens, parseIdentifierToken(cursor))
			continue
		}
		if isNumeric(char) {
			tokens = append(tokens, parseInt(cursor))
			continue
		}
		fmt.Printf("unexpected character: %c\n", char)
		break
	}
	return tokens
}
