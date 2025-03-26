package parser

import (
	"fmt"
	"unicode"
)

type tokenizer struct {
	chars []rune
	loc   Location
}

func (t *tokenizer) tokenize() []Token {
	var tokens []Token
	for {
		if len(t.chars) == 0 {
			break
		}
		char := t.char()
		if isIdentifierFirstChar(char) {
			tokens = append(tokens, t.identifier())
			continue
		}
		if unicode.IsSpace(char) {
			t.skipWhitespace()
			continue
		}
		if char == '"' || char == '\'' {
			literal, err := t.stringLiteral()
			if err != nil {
				break
			}
			tokens = append(tokens, literal)
			continue
		}
		if (char >= '0' && char <= '9') || char == '-' {
			literal, err := t.intLiteral()
			if err != nil {
				break
			}
			tokens = append(tokens, literal)
			continue
		}
		if char == ':' {
			start := t.loc
			t.next()
			if t.char() == ':' {
				tokens = append(tokens, NewSymbolToken(DoubleColon, NewSpan(start, t.loc)))
				t.next()
				continue
			}
			tokens = append(tokens, NewSymbolToken(Colon, NewSpan(start, start)))
			continue
		}
		span := NewSpan(t.loc, t.loc)
		switch char {
		case '<':
			tokens = append(tokens, NewSymbolToken(Lt, span))
		case '>':
			tokens = append(tokens, NewSymbolToken(Gt, span))
		case ',':
			tokens = append(tokens, NewSymbolToken(Comma, span))
		case '{':
			tokens = append(tokens, NewSymbolToken(Lbrace, span))
		case '}':
			tokens = append(tokens, NewSymbolToken(Rbrace, span))
		case '|':
			tokens = append(tokens, NewSymbolToken(Pipe, span))
		case '&':
			tokens = append(tokens, NewSymbolToken(Amp, span))
		case '(':
			tokens = append(tokens, NewSymbolToken(Lparen, span))
		case ')':
			tokens = append(tokens, NewSymbolToken(Rparen, span))
		case '=':
			tokens = append(tokens, NewSymbolToken(Eq, span))
		case '*':
			tokens = append(tokens, NewSymbolToken(Asterisk, span))
		}
		t.next()
	}
	return tokens
}

func (t *tokenizer) char() rune {
	if len(t.chars) == 0 {
		return 0
	}
	return t.chars[0]
}

func (t *tokenizer) identifier() Token {
	var name []rune
	start := t.loc
	end := t.loc
	for {
		if !isIdentifierChar(t.char()) {
			break
		}
		name = append(name, t.char())
		end = t.loc
		t.next()
	}
	return Token{Kind: Identifier, Val: string(name), Loc: NewSpan(start, end)}
}

func isIdentifierFirstChar(char rune) bool {
	return char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z'
}

func isIdentifierChar(char rune) bool {
	return isIdentifierFirstChar(char) || char == '-' || char == '*' || char == '_'
}

func isWhitespace(char rune) bool {
	return unicode.IsSpace(char)
}

func (t *tokenizer) next() {
	if t.char() == '\n' {
		t.loc.Line++
		t.loc.Col = 1
	} else {
		t.loc.Col++
	}
	t.chars = t.chars[1:]
}

func (t *tokenizer) skipWhitespace() {
	for {
		if !isWhitespace(t.char()) {
			break
		}
		t.next()
	}
}

func (t *tokenizer) stringLiteral() (Token, error) {
	quote := t.char()
	start := t.loc
	t.next()
	var str []rune
	for {
		char := t.char()
		if char == 0 || char == '\n' {
			return Token{}, fmt.Errorf("unterminated string literal")
		}
		if char == quote {
			end := t.loc
			t.next()
			return Token{Kind: StringLiteral, Val: string(str), Loc: NewSpan(start, end)}, nil
		}
		str = append(str, char)
		t.next()
	}
}

func (t *tokenizer) intLiteral() (Token, error) {
	start := t.loc
	end := t.loc
	if t.char() == '0' {
		t.next()
		if t.char() >= '0' && t.char() <= '9' {
			return Token{}, fmt.Errorf("integer literal cannot have leading zero")
		}
		return Token{Kind: IntLiteral, Val: "0", Loc: NewSpan(start, end)}, nil
	}
	var chars []rune
	if t.char() == '-' {
		chars = append(chars, '-')
		t.next()
		char := t.char()
		if char < '1' || char > '9' {
			return Token{}, fmt.Errorf("invalid character %c after '-'. Expected 1-9", char)
		}
	}
	for {
		char := t.char()
		if char < '0' || char > '9' {
			break
		}
		chars = append(chars, char)
		end = t.loc
		t.next()
	}
	return Token{Kind: IntLiteral, Val: string(chars), Loc: NewSpan(start, end)}, nil
}

func newTokenizer(src string) *tokenizer {
	return &tokenizer{
		chars: []rune(src),
		loc:   Location{Line: 1, Col: 1},
	}
}

func Tokenize(src string) []Token {
	return newTokenizer(src).tokenize()
}
