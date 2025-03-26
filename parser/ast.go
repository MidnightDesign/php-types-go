package parser

import (
	"fmt"
	"strings"
)

type Node interface {
	fmt.Stringer
}

type IdentifierNode struct {
	Name          string
	TypeArguments []Node
}

type CurlyListNode struct {
	Name     string
	Elements []Node
}
type CurlyKeyValueNode struct {
	Name    string
	Members []*MemberNode
}

type MemberNode struct {
	Key      string
	Value    Node
	Optional bool
}

func NewMember(key string, value Node) *MemberNode {
	return &MemberNode{Key: key, Value: value, Optional: false}
}

func NewOptionalMember(key string, value Node) *MemberNode {
	return &MemberNode{Key: key, Value: value, Optional: true}
}

func (n MemberNode) String() string {
	if n.Optional {
		return fmt.Sprintf("%s?: %s", n.Key, n.Value)
	}
	return fmt.Sprintf("%s: %s", n.Key, n.Value)
}

type CallableNode struct {
	ReturnType Node
	Parameters []*ParamNode
}

type ParamNode struct {
	Type     Node
	Optional bool
}

func NewParam(typeNode Node) *ParamNode {
	return &ParamNode{Type: typeNode, Optional: false}
}

func NewOptionalParam(typeNode Node) *ParamNode {
	return &ParamNode{Type: typeNode, Optional: true}
}

func (n *ParamNode) String() string {
	if n.Optional {
		return fmt.Sprintf("%s=", n.Type.String())
	}
	return n.Type.String()
}

type StringLiteralNode struct {
	Value string
}

type IntLiteralNode struct {
	Value int
}

type UnionNode struct {
	Elements []Node
}

type IntersectionNode struct {
	Elements []Node
}

func nodeList(list []Node) string {
	if len(list) == 0 {
		return ""
	}
	elements := make([]string, len(list))
	for i, element := range list {
		elements[i] = element.String()
	}
	return strings.Join(elements, ", ")
}

func (n IdentifierNode) String() string {
	if len(n.TypeArguments) == 0 {
		return n.Name
	}
	return fmt.Sprintf("%s<%s>", n.Name, nodeList(n.TypeArguments))
}

func (n CurlyListNode) String() string {
	return fmt.Sprintf("%s{%s}", n.Name, nodeList(n.Elements))
}

func (n CurlyKeyValueNode) String() string {
	members := make([]string, len(n.Members))
	for i, member := range n.Members {
		members[i] = member.String()
	}
	return fmt.Sprintf("%s{%s}", n.Name, strings.Join(members, ", "))
}

func (n CallableNode) String() string {
	parameters := make([]string, len(n.Parameters))
	for i, parameter := range n.Parameters {
		parameters[i] = parameter.String()
	}
	return fmt.Sprintf("callable(%s): %s", strings.Join(parameters, ", "), n.ReturnType)
}

func (n StringLiteralNode) String() string {
	return fmt.Sprintf("\"%s\"", n.Value)
}

func (n IntLiteralNode) String() string {
	return fmt.Sprintf("%d", n.Value)
}

func (n UnionNode) String() string {
	elements := make([]string, len(n.Elements))
	for i, element := range n.Elements {
		elements[i] = element.String()
	}
	return strings.Join(elements, " | ")
}

func (n IntersectionNode) String() string {
	elements := make([]string, len(n.Elements))
	for i, element := range n.Elements {
		elements[i] = element.String()
	}
	return strings.Join(elements, " & ")
}

func NewSimpleNode(name string) Node {
	return &IdentifierNode{Name: name}
}

func NewGenericNode(name string, typeArguments []Node) Node {
	return &IdentifierNode{Name: name, TypeArguments: typeArguments}
}

func NewCurlyListNode(name string, elements []Node) Node {
	return &CurlyListNode{Name: name, Elements: elements}
}

func NewCurlyKeyValueNode(name string, members []*MemberNode) Node {
	return &CurlyKeyValueNode{Name: name, Members: members}
}

func NewCallableNode(returnType Node, parameters []*ParamNode) Node {
	return &CallableNode{ReturnType: returnType, Parameters: parameters}
}

func NewStringLiteralNode(value string) Node {
	return &StringLiteralNode{Value: value}
}

func NewIntLiteralNode(value int) Node {
	return &IntLiteralNode{Value: value}
}

func NewUnionNode(first Node, second Node, other ...Node) Node {
	elements := append([]Node{first, second}, other...)
	return &UnionNode{Elements: elements}
}

func NewIntersectionNode(first Node, second Node, other ...Node) Node {
	elements := append([]Node{first, second}, other...)
	return &IntersectionNode{Elements: elements}
}
