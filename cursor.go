package main

type cursor[T comparable] struct {
	items   []T
	current int
}

func (c *cursor[T]) peek() (T, bool) {
	if c.current >= len(c.items) {
		return *new(T), true
	}
	return c.items[c.current], false
}

func (c *cursor[T]) next() {
	c.current++
}

func (c *cursor[T]) isCurrent(v T) bool {
	item, eof := c.peek()
	if eof {
		return false
	}
	return item == v
}

func (c *cursor[T]) currentMatches(v func(T) bool) bool {
	item, eof := c.peek()
	if eof {
		return false
	}
	return v(item)
}

func newCursor[T comparable](s []T) *cursor[T] {
	return &cursor[T]{s, 0}
}
