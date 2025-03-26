package parser

import "fmt"

type Location struct {
	Line int
	Col  int
}

func (l Location) String() string {
	return fmt.Sprintf("%d:%d", l.Line, l.Col)
}

func (l Location) add(i int) Location {
	return Location{Line: l.Line, Col: l.Col + i}
}

func NewLocation(line, col int) Location {
	return Location{Line: line, Col: col}
}

type Span struct {
	Start Location
	End   Location
}

func (s Span) String() string {
	return fmt.Sprintf("%s-%s", s.Start, s.End)
}

func NewSpan(start, end Location) Span {
	return Span{Start: start, End: end}
}

func NewSpanFromInts(startLine, startCol, endLine, endCol int) Span {
	return Span{Start: NewLocation(startLine, startCol), End: NewLocation(endLine, endCol)}
}

func NewSingleCharSpan(line, col int) Span {
	return NewSpanFromInts(line, col, line, col)
}
