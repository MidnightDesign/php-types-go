package parser

import (
	"fmt"
	"strconv"
)

type tokenKind int8

func (t tokenKind) String() string {
	switch t {
	case Identifier:
		return "Identifier"
	case StringLiteral:
		return "StringLiteral"
	case IntLiteral:
		return "IntLiteral"
	case Gt:
		return ">"
	case Lt:
		return "<"
	case Comma:
		return ","
	case Lbrace:
		return "{"
	case Rbrace:
		return "}"
	case Pipe:
		return "|"
	case Amp:
		return "&"
	case Colon:
		return ":"
	case Lparen:
		return "("
	case Rparen:
		return ")"
	case Eq:
		return "="
	case DoubleColon:
		return "::"
	case Asterisk:
		return "*"
	}
	return "unknown"
}

const (
	Identifier tokenKind = iota
	Gt
	Lt
	Comma
	Lbrace
	Rbrace
	StringLiteral
	Pipe
	Amp
	Colon
	IntLiteral
	Lparen
	Rparen
	Eq
	DoubleColon
	Asterisk
)

type Token struct {
	Kind tokenKind
	Val  string
	Loc  Span
}

func (t Token) String() string {
	var v string
	switch t.Kind {
	case Identifier:
		v = t.Val
	case StringLiteral:
		v = fmt.Sprintf("\"%s\"", t.Val)
	case IntLiteral:
		v = fmt.Sprintf("%s", t.Val)
	default:
		v = t.Kind.String()
	}
	return fmt.Sprintf("%s (%s)", v, t.Loc)
}

func NewIdentifierToken(name string, loc Span) Token {
	return Token{Kind: Identifier, Val: name, Loc: loc}
}

func NewStringLiteralToken(name string, loc Span) Token {
	return Token{Kind: StringLiteral, Val: name, Loc: loc}
}

func NewIntLiteralToken(value int, loc Span) Token {
	return Token{Kind: IntLiteral, Val: strconv.Itoa(value), Loc: loc}
}

func NewSymbolToken(kind tokenKind, loc Span) Token {
	return Token{Kind: kind, Loc: loc}
}
