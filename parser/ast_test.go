package parser_test

import (
	"github.com/MidnightDesign/php-types-go/parser"
	"testing"
)

func TestNode_String(t *testing.T) {
	var tests = []struct {
		name string
		node parser.Node
		want string
	}{
		{
			node: parser.NewSimpleNode("string"),
			want: "string",
		},
		{
			node: parser.NewGenericNode("list", []parser.Node{parser.NewSimpleNode("int")}),
			want: "list<int>",
		},
		{
			node: parser.NewGenericNode("array", []parser.Node{
				parser.NewSimpleNode("array-key"),
				parser.NewSimpleNode("string"),
			}),
			want: "array<array-key, string>",
		},
		{
			node: parser.NewCurlyListNode("array", nil),
			want: "array{}",
		},
		{
			node: parser.NewCurlyListNode("array", []parser.Node{
				parser.NewSimpleNode("string"),
			}),
			want: "array{string}",
		},
		{
			node: parser.NewCurlyListNode("array", []parser.Node{
				parser.NewSimpleNode("string"),
				parser.NewSimpleNode("int"),
			}),
			want: "array{string, int}",
		},
		{
			node: parser.NewCurlyKeyValueNode("array", nil),
			want: "array{}",
		},
		{
			node: parser.NewCurlyKeyValueNode("array", []*parser.MemberNode{
				parser.NewMember("foo", parser.NewSimpleNode("string")),
			}),
			want: "array{foo: string}",
		},
		{
			node: parser.NewCurlyKeyValueNode("array", []*parser.MemberNode{
				parser.NewMember("foo", parser.NewSimpleNode("string")),
				parser.NewMember("bar", parser.NewSimpleNode("int")),
			}),
			want: "array{foo: string, bar: int}",
		},
		{
			node: parser.NewCurlyKeyValueNode("object", []*parser.MemberNode{
				parser.NewOptionalMember("foo", parser.NewSimpleNode("string")),
				parser.NewMember("bar", parser.NewSimpleNode("int")),
			}),
			want: "object{foo?: string, bar: int}",
		},
		{
			node: parser.NewCallableNode(parser.NewSimpleNode("void"), nil),
			want: "callable(): void",
		},
		{
			node: parser.NewCallableNode(parser.NewSimpleNode("bool"), []*parser.ParamNode{
				parser.NewParam(parser.NewSimpleNode("string")),
				parser.NewParam(parser.NewSimpleNode("int")),
			}),
			want: "callable(string, int): bool",
		},
		{
			node: parser.NewCallableNode(parser.NewSimpleNode("bool"), []*parser.ParamNode{
				parser.NewParam(parser.NewSimpleNode("string")),
				parser.NewOptionalParam(parser.NewSimpleNode("int")),
			}),
			want: "callable(string, int=): bool",
		},
		{
			node: parser.NewStringLiteralNode(""),
			want: "\"\"",
		},
		{
			node: parser.NewStringLiteralNode("foo"),
			want: "\"foo\"",
		},
		{
			node: parser.NewIntLiteralNode(-23),
			want: "-23",
		},
		{
			node: parser.NewIntLiteralNode(0),
			want: "0",
		},
		{
			node: parser.NewIntLiteralNode(42),
			want: "42",
		},
		{
			node: parser.NewUnionNode(parser.NewSimpleNode("string"), parser.NewSimpleNode("int")),
			want: "string | int",
		},
		{
			node: parser.NewIntersectionNode(
				parser.NewCurlyKeyValueNode("array", []*parser.MemberNode{
					parser.NewMember("foo", parser.NewSimpleNode("string")),
				}),
				parser.NewCurlyKeyValueNode("array", []*parser.MemberNode{
					parser.NewMember("bar", parser.NewSimpleNode("int")),
				}),
			),
			want: "array{foo: string} & array{bar: int}",
		},
	}

	for _, test := range tests {
		name := test.name
		if name == "" {
			name = test.want
		}
		t.Run(name, func(t *testing.T) {
			if got := test.node.String(); got != test.want {
				t.Errorf("Node.String() = %v, want %v", got, test.want)
			}
		})
	}
}
