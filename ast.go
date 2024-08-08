package main

import "fmt"

type kind int

const (
	simple kind = iota + 1
	angle
	curly
	keyValue
	optional
	integerNode
	stringNode
	ellipsisNode
	union
	intersection
	callableNode
	parameters
	parameter
)

func createOptionalNode() *node {
	return &node{kind: optional}
}

func createParameterNode(typeNode *node, optional bool) node {
	children := []node{*typeNode}
	if optional {
		children = append(children, *createOptionalNode())
	}
	return node{kind: parameter, children: children}
}

func createParametersNode(params []node) node {
	return node{kind: parameters, children: params}
}

func createCurlyNode(name string, items []node) *node {
	return &node{kind: curly, value: name, children: items}
}

func createKeyValueNode(key string, value node, optional bool) *node {
	children := []node{value}
	if optional {
		children = append(children, *createOptionalNode())
	}
	return &node{kind: keyValue, value: key, children: children}
}

func createAngleNode(name string, args []node) *node {
	return &node{kind: angle, value: name, children: args}
}

type node struct {
	kind     kind
	value    string
	children []node
}

// array<string, int>
//
//	============= Parsed by this function.
//
// ===== The name of the type is passed in as value.
func parseAngle(c *cursor[token], value string) *node {
	arguments := make([]node, 0)
	c.next()
	for {
		argument := parseNode(c)
		if argument == nil {
			break
		}
		arguments = append(arguments, *argument)
		if currentIs(comma, c) {
			c.next()
			continue
		}
		break
	}
	if !currentIs(closeAngle, c) {
		return nil
	}
	c.next()
	return createAngleNode(value, arguments)
}

func currentIs(tokenType tokenType, c *cursor[token]) bool {
	return c.currentMatches(func(t token) bool { return t.t == tokenType })
}

func parseCurlyItem(c *cursor[token]) *node {
	n := parseNode(c)
	if n == nil {
		return nil
	}
	isColonOrQuestion := c.currentMatches(func(t token) bool { return t.t == colon || t.t == question })
	if !isColonOrQuestion {
		return n
	}
	opt := false
	if currentIs(question, c) {
		opt = true
		c.next()
	}
	if !currentIs(colon, c) {
		return n
	}
	c.next()
	return createKeyValueNode(n.value, *parseNode(c), opt)
}

func parseCurly(c *cursor[token], name string) *node {
	items := make([]node, 0)
	c.next()
	t, eof := c.peek()
	if eof {
		return nil
	}
	if t.t == closeCurly {
		c.next()
		return &node{kind: curly, value: name}
	}
	for {
		item := parseCurlyItem(c)
		if item == nil {
			break
		}
		items = append(items, *item)
		t, eof := c.peek()
		if eof {
			break
		}
		if t.t == comma {
			c.next()
			continue
		}
		break
	}
	t, eof = c.peek()
	if eof {
		return nil
	}
	if t.t != closeCurly {
		return nil
	}
	c.next()
	return createCurlyNode(name, items)
}

func parseCallable(tokens *cursor[token], name string) *node {
	tokens.next()
	params := make([]node, 0)
	for {
		t, eof := tokens.peek()
		if eof {
			break
		}
		if t.t == closeBrace {
			break
		}
		parameter := parseNode(tokens)
		if parameter == nil {
			break
		}
		optional := false
		if tokens.currentMatches(func(t token) bool { return t.t == equal }) {
			tokens.next()
			optional = true
		}
		params = append(params, createParameterNode(parameter, optional))
		t, eof = tokens.peek()
		if eof {
			break
		}
		if t.t == comma {
			tokens.next()
			continue
		}
		break
	}
	token, eof := tokens.peek()
	if eof {
		return nil
	}
	if token.t != closeBrace {
		return nil
	}
	tokens.next()
	token, eof = tokens.peek()
	if eof {
		return &node{kind: callableNode, value: name, children: params}
	}
	if token.t != colon {
		return &node{kind: callableNode, value: name, children: params}
	}
	tokens.next()
	callableChildren := []node{createParametersNode(params)}
	returnType := parseNode(tokens)
	if returnType != nil {
		callableChildren = append(callableChildren, *returnType)
	}
	return &node{kind: callableNode, value: name, children: callableChildren}
}

func parseIdentifier(c *cursor[token]) *node {
	identifierToken, eof := c.peek()
	if eof {
		return nil
	}
	c.next()
	t, eof := c.peek()
	if eof {
		return &node{kind: simple, value: identifierToken.stringValue}
	}
	if t.t == openBrace {
		return parseCallable(c, identifierToken.stringValue)
	}
	if t.t == openAngle {
		return parseAngle(c, identifierToken.stringValue)
	}
	if t.t == openCurly {
		return parseCurly(c, identifierToken.stringValue)
	}
	return &node{kind: simple, value: identifierToken.stringValue}
}

func parseNegativeIntLiteral(tokens *cursor[token]) *node {
	tokens.next()
	t, eof := tokens.peek()
	if eof {
		return nil
	}
	if t.t != integer {
		return nil
	}
	tokens.next()
	return &node{kind: integerNode, value: "-" + t.stringValue}
}

func parseNode(tokens *cursor[token]) *node {
	var left *node
	for {
		t, eof := tokens.peek()
		if eof {
			break
		}
		switch t.t {
		case identifier:
			left = parseIdentifier(tokens)
			continue
		case hyphen:
			left = parseNegativeIntLiteral(tokens)
			continue
		case integer:
			tokens.next()
			left = &node{kind: integerNode, value: t.stringValue}
			continue
		case stringLiteral:
			tokens.next()
			left = &node{kind: stringNode, value: t.stringValue}
			continue
		case ellipsis:
			tokens.next()
			left = &node{kind: ellipsisNode}
			continue
		case pipe:
			tokens.next()
			if left == nil {
				fmt.Printf("Unexpected token: %v\n", t.t)
				panic("Unexpected token")
			}
			left = &node{kind: union, children: []node{*left, *parseNode(tokens)}}
			continue
		case ampersand:
			tokens.next()
			if left == nil {
				fmt.Printf("Unexpected token: %v\n", t.t)
				panic("Unexpected token")
			}
			left = &node{kind: intersection, children: []node{*left, *parseNode(tokens)}}
			continue
		}
		if left != nil {
			return left
		}
		fmt.Printf("Unexpected token: %v\n", t.t)
		panic("Unexpected token")
	}
	return left
}

func (n *node) print() string {
	if n.kind == simple || n.kind == integerNode {
		return n.value
	}
	if n.kind == stringNode {
		return "'" + n.value + "'"
	}
	if n.kind == ellipsisNode {
		return "..."
	}
	if n.kind == angle {
		args := ""
		for i, arg := range n.children {
			if i > 0 {
				args += ", "
			}
			args += arg.print()
		}
		return n.value + "<" + args + ">"
	}
	if n.kind == curly {
		items := ""
		for i, item := range n.children {
			if i > 0 {
				items += ", "
			}
			items += item.print()
		}
		return n.value + "{" + items + "}"
	}
	if n.kind == union {
		return n.children[0].print() + " | " + n.children[1].print()
	}
	if n.kind == intersection {
		return n.children[0].print() + " & " + n.children[1].print()
	}
	if n.kind == callableNode {
		params := ""
		for i, param := range n.children[0].children {
			if i > 0 {
				params += ", "
			}
			params += param.print()
		}
		if len(n.children) == 1 {
			return n.value + "(" + params + ")"
		}
		return n.value + "(" + params + "): " + n.children[1].print()
	}
	if n.kind == parameter {
		if len(n.children) == 1 {
			return n.children[0].print()
		}
		return n.children[0].print() + "="
	}
	if n.kind == keyValue {
		if len(n.children) == 1 {
			return n.value + ": " + n.children[0].print()
		}
		return n.value + "?: " + n.children[0].print()
	}
	return ""
}
