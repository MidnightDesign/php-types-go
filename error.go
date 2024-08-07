package main

import "fmt"

type syntaxError struct {
	message  string
	location span
}

func newSyntaxError(message string, location span) syntaxError {
	return syntaxError{message, location}
}

func (e syntaxError) Error() string {
	return fmt.Sprintf("%s at %s", e.message, e.location)
}
