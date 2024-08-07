package main

import "fmt"

type location struct {
	line   int
	column int
}

func newLocation(line int, column int) location {
	return location{line, column}
}

func (l location) String() string {
	return fmt.Sprintf("%d:%d", l.line, l.column)
}

type span struct {
	start location
	end   location
}

func newSpan(start location, end location) span {
	return span{start, end}
}

func (s span) String() string {
	return s.start.String() + "-" + s.end.String()
}

func newSingleCharSpan(l location) span {
	return span{start: l, end: l}
}
