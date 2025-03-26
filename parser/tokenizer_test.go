package parser_test

import (
	"github.com/MidnightDesign/php-types-go/parser"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		src    string
		tokens []parser.Token
	}{
		{"string", []parser.Token{parser.NewIdentifierToken("string", parser.NewSpanFromInts(1, 1, 1, 6))}},
		{"list<string>", []parser.Token{
			parser.NewIdentifierToken("list", parser.NewSpanFromInts(1, 1, 1, 4)),
			parser.NewSymbolToken(parser.Lt, parser.NewSingleCharSpan(1, 5)),
			parser.NewIdentifierToken("string", parser.NewSpanFromInts(1, 6, 1, 11)),
			parser.NewSymbolToken(parser.Gt, parser.NewSingleCharSpan(1, 12)),
		}},
		{"array<string, string>", []parser.Token{
			parser.NewIdentifierToken("array", parser.NewSpanFromInts(1, 1, 1, 5)),
			parser.NewSymbolToken(parser.Lt, parser.NewSingleCharSpan(1, 6)),
			parser.NewIdentifierToken("string", parser.NewSpanFromInts(1, 7, 1, 12)),
			parser.NewSymbolToken(parser.Comma, parser.NewSingleCharSpan(1, 13)),
			parser.NewIdentifierToken("string", parser.NewSpanFromInts(1, 15, 1, 20)),
			parser.NewSymbolToken(parser.Gt, parser.NewSingleCharSpan(1, 21)),
		}},
		{"array-key", []parser.Token{parser.NewIdentifierToken("array-key", parser.NewSpanFromInts(1, 1, 1, 9))}},
		{"array{}", []parser.Token{
			parser.NewIdentifierToken("array", parser.NewSpanFromInts(1, 1, 1, 5)),
			parser.NewSymbolToken(parser.Lbrace, parser.NewSingleCharSpan(1, 6)),
			parser.NewSymbolToken(parser.Rbrace, parser.NewSingleCharSpan(1, 7)),
		}},
		{"array{int}", []parser.Token{
			parser.NewIdentifierToken("array", parser.NewSpanFromInts(1, 1, 1, 5)),
			parser.NewSymbolToken(parser.Lbrace, parser.NewSingleCharSpan(1, 6)),
			parser.NewIdentifierToken("int", parser.NewSpanFromInts(1, 7, 1, 9)),
			parser.NewSymbolToken(parser.Rbrace, parser.NewSingleCharSpan(1, 10)),
		}},
		{"array{int, string}", []parser.Token{
			parser.NewIdentifierToken("array", parser.NewSpanFromInts(1, 1, 1, 5)),
			parser.NewSymbolToken(parser.Lbrace, parser.NewSingleCharSpan(1, 6)),
			parser.NewIdentifierToken("int", parser.NewSpanFromInts(1, 7, 1, 9)),
			parser.NewSymbolToken(parser.Comma, parser.NewSingleCharSpan(1, 10)),
			parser.NewIdentifierToken("string", parser.NewSpanFromInts(1, 12, 1, 17)),
			parser.NewSymbolToken(parser.Rbrace, parser.NewSingleCharSpan(1, 18)),
		}},
		{"\"\"", []parser.Token{parser.NewStringLiteralToken("", parser.NewSpanFromInts(1, 1, 1, 2))}},
		{"''", []parser.Token{parser.NewStringLiteralToken("", parser.NewSpanFromInts(1, 1, 1, 2))}},
		{"\"foo\"", []parser.Token{parser.NewStringLiteralToken("foo", parser.NewSpanFromInts(1, 1, 1, 5))}},
		{"'foo'", []parser.Token{parser.NewStringLiteralToken("foo", parser.NewSpanFromInts(1, 1, 1, 5))}},
		{"string | int", []parser.Token{
			parser.NewIdentifierToken("string", parser.NewSpanFromInts(1, 1, 1, 6)),
			parser.NewSymbolToken(parser.Pipe, parser.NewSingleCharSpan(1, 8)),
			parser.NewIdentifierToken("int", parser.NewSpanFromInts(1, 10, 1, 12)),
		}},
		{"Foo & Bar", []parser.Token{
			parser.NewIdentifierToken("Foo", parser.NewSpanFromInts(1, 1, 1, 3)),
			parser.NewSymbolToken(parser.Amp, parser.NewSingleCharSpan(1, 5)),
			parser.NewIdentifierToken("Bar", parser.NewSpanFromInts(1, 7, 1, 9)),
		}},
		{"array{foo: bool}", []parser.Token{
			parser.NewIdentifierToken("array", parser.NewSpanFromInts(1, 1, 1, 5)),
			parser.NewSymbolToken(parser.Lbrace, parser.NewSingleCharSpan(1, 6)),
			parser.NewIdentifierToken("foo", parser.NewSpanFromInts(1, 7, 1, 9)),
			parser.NewSymbolToken(parser.Colon, parser.NewSingleCharSpan(1, 10)),
			parser.NewIdentifierToken("bool", parser.NewSpanFromInts(1, 12, 1, 15)),
			parser.NewSymbolToken(parser.Rbrace, parser.NewSingleCharSpan(1, 16)),
		}},
		{"420", []parser.Token{parser.NewIntLiteralToken(420, parser.NewSpanFromInts(1, 1, 1, 3))}},
		{"0", []parser.Token{parser.NewIntLiteralToken(0, parser.NewSpanFromInts(1, 1, 1, 1))}},
		{"-420", []parser.Token{parser.NewIntLiteralToken(-420, parser.NewSpanFromInts(1, 1, 1, 4))}},
		{"callable(string, int=): void", []parser.Token{
			parser.NewIdentifierToken("callable", parser.NewSpanFromInts(1, 1, 1, 8)),
			parser.NewSymbolToken(parser.Lparen, parser.NewSingleCharSpan(1, 9)),
			parser.NewIdentifierToken("string", parser.NewSpanFromInts(1, 10, 1, 15)),
			parser.NewSymbolToken(parser.Comma, parser.NewSingleCharSpan(1, 16)),
			parser.NewIdentifierToken("int", parser.NewSpanFromInts(1, 18, 1, 20)),
			parser.NewSymbolToken(parser.Eq, parser.NewSingleCharSpan(1, 21)),
			parser.NewSymbolToken(parser.Rparen, parser.NewSingleCharSpan(1, 22)),
			parser.NewSymbolToken(parser.Colon, parser.NewSingleCharSpan(1, 23)),
			parser.NewIdentifierToken("void", parser.NewSpanFromInts(1, 25, 1, 28)),
		}},
		{"Foo::BAR", []parser.Token{
			parser.NewIdentifierToken("Foo", parser.NewSpanFromInts(1, 1, 1, 3)),
			parser.NewSymbolToken(parser.DoubleColon, parser.NewSpanFromInts(1, 4, 1, 5)),
			parser.NewIdentifierToken("BAR", parser.NewSpanFromInts(1, 6, 1, 8)),
		}},
		{"Foo::*", []parser.Token{
			parser.NewIdentifierToken("Foo", parser.NewSpanFromInts(1, 1, 1, 3)),
			parser.NewSymbolToken(parser.DoubleColon, parser.NewSpanFromInts(1, 4, 1, 5)),
			parser.NewSymbolToken(parser.Asterisk, parser.NewSingleCharSpan(1, 6)),
		}},
		{"Foo::STATUS_*", []parser.Token{
			parser.NewIdentifierToken("Foo", parser.NewSpanFromInts(1, 1, 1, 3)),
			parser.NewSymbolToken(parser.DoubleColon, parser.NewSpanFromInts(1, 4, 1, 5)),
			parser.NewIdentifierToken("STATUS_*", parser.NewSpanFromInts(1, 6, 1, 13)),
		}},
		{"array{\n    foo: string,\n}", []parser.Token{
			parser.NewIdentifierToken("array", parser.NewSpanFromInts(1, 1, 1, 5)),
			parser.NewSymbolToken(parser.Lbrace, parser.NewSingleCharSpan(1, 6)),
			parser.NewIdentifierToken("foo", parser.NewSpanFromInts(2, 5, 2, 7)),
			parser.NewSymbolToken(parser.Colon, parser.NewSingleCharSpan(2, 8)),
			parser.NewIdentifierToken("string", parser.NewSpanFromInts(2, 10, 2, 15)),
			parser.NewSymbolToken(parser.Comma, parser.NewSingleCharSpan(2, 16)),
			parser.NewSymbolToken(parser.Rbrace, parser.NewSingleCharSpan(3, 1)),
		}},
		{"string | \"foo", []parser.Token{
			parser.NewIdentifierToken("string", parser.NewSpanFromInts(1, 1, 1, 6)),
			parser.NewSymbolToken(parser.Pipe, parser.NewSingleCharSpan(1, 8)),
		}},
		{"string | 023", []parser.Token{
			parser.NewIdentifierToken("string", parser.NewSpanFromInts(1, 1, 1, 6)),
			parser.NewSymbolToken(parser.Pipe, parser.NewSingleCharSpan(1, 8)),
		}},
		{"string | -0", []parser.Token{
			parser.NewIdentifierToken("string", parser.NewSpanFromInts(1, 1, 1, 6)),
			parser.NewSymbolToken(parser.Pipe, parser.NewSingleCharSpan(1, 8)),
		}},
	}
	for _, test := range tests {
		t.Run(test.src, func(t *testing.T) {
			tokens := parser.Tokenize(test.src)
			if len(tokens) != len(test.tokens) {
				t.Errorf("\"%s\": expected %d tokens, got %d", test.src, len(test.tokens), len(tokens))
				return
			}
			for i, actual := range tokens {
				if actual != test.tokens[i] {
					t.Errorf("expected token %v, got %v", test.tokens[i], actual)
				}
			}
		})
	}
}
