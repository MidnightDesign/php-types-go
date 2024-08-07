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

func (l *lexer) parseIdentifierToken() token {
	char, _ := l.chars.peek()
	idStr := string(char)
	l.next()
	for {
		char, eof := l.chars.peek()
		if eof {
			break
		}
		if !isIdentifierChar(char) {
			break
		}
		idStr += string(char)
		l.next()
	}
	return token{t: identifier, stringValue: idStr}
}

func (l *lexer) parseInt() token {
	char, _ := l.chars.peek()
	numStr := string(char)
	l.next()
	for {
		char, eof := l.chars.peek()
		if eof {
			break
		}
		if !isNumeric(char) {
			break
		}
		numStr += string(char)
		l.next()
	}
	return token{t: integer, stringValue: numStr}
}

func (l *lexer) skipWhitespace() {
	for {
		char, eof := l.chars.peek()
		if eof {
			break
		}
		if !unicode.IsSpace(char) {
			break
		}
		l.next()
	}
}

func (l *lexer) parseStringLiteralToken() (token, error) {
	l.next()
	str := ""
	for {
		char, eof := l.chars.peek()
		if eof {
			return token{}, newSyntaxError("unexpected end of input", newSingleCharSpan(l.location))
		}
		if char == '\n' {
			return token{}, newSyntaxError("unexpected newline in string literal", newSingleCharSpan(l.location))
		}
		if char == '\'' {
			l.next()
			break
		}
		str += string(char)
		l.next()
	}
	return token{t: stringLiteral, stringValue: str}, nil
}

func (l *lexer) parseEllipsis() (token, error) {
	periods := 0
	for {
		char, eof := l.chars.peek()
		if eof {
			return token{}, newSyntaxError("unexpected end of input", newSingleCharSpan(l.location))
		}
		if char != '.' {
			return token{}, newSyntaxError(
				fmt.Sprintf("unexpected character: %c", char),
				newSingleCharSpan(l.location),
			)
		}
		periods++
		l.next()
		if periods != 3 {
			continue
		}
		break
	}
	return token{t: ellipsis, stringValue: "..."}, nil
}

func (l *lexer) next() {
	char, eof := l.chars.peek()
	if eof {
		return
	}
	l.chars.next()
	if char == '\n' {
		l.location.column++
		return
	}
	l.location.line++
	l.location.column = 1
}

func (l *lexer) lex() ([]token, error) {
	tokens := make([]token, 0)
	for {
		char, eof := l.chars.peek()
		if eof {
			break
		}
		tokenType := singleCharTokens[char]
		if tokenType != 0 {
			tokens = append(tokens, token{t: tokenType, stringValue: string(char)})
			l.next()
			continue
		}
		if char == '.' {
			t, err := l.parseEllipsis()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, t)
			continue
		}
		if char == '\'' {
			t, err := l.parseStringLiteralToken()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, t)
			continue
		}
		if unicode.IsSpace(char) {
			l.skipWhitespace()
			continue
		}
		if isAlpha(char) {
			tokens = append(tokens, l.parseIdentifierToken())
			continue
		}
		if isNumeric(char) {
			tokens = append(tokens, l.parseInt())
			continue
		}
		return nil, newSyntaxError(
			fmt.Sprintf("unexpected character: %c", char),
			newSingleCharSpan(l.location),
		)
	}
	return tokens, nil
}

type lexer struct {
	chars    *cursor[rune]
	location location
}

func newLexer(s string) *lexer {
	chars := newCursor([]rune(s))
	return &lexer{
		chars:    chars,
		location: newLocation(1, 1),
	}
}
