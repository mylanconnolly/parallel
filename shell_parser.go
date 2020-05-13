package main

import (
	"strings"
	"unicode"
)

// This function is meant to split an input string into shell words. For example
//
//     "one two" three 'four five six'
//
// would become
//
//     []string{"one two", "three", "four five six"}
//
// I'm doing this so that I don't have to resort to hacky solutions like running
// a shell or smoething.
func shellParser(input string) []string {
	out := []string{}

	var (
		buf         strings.Builder
		escape      bool
		doubleQuote bool
		singleQuote bool
		gotWord     bool
	)
	for _, r := range input {
		switch {
		case escape:
			buf.WriteRune(r)
			escape = false
			continue
		case unicode.IsSpace(r):
			if singleQuote || doubleQuote {
				buf.WriteRune(r)
			} else if gotWord {
				out = append(out, buf.String())
				buf.Reset()
				gotWord = false
			}
			continue
		case r == '\\':
			if singleQuote {
				buf.WriteRune(r)
			} else {
				escape = true
			}
			continue
		case r == '"':
			if !singleQuote {
				if doubleQuote {
					gotWord = true
				}
				doubleQuote = !doubleQuote
				continue
			}
		case r == '\'':
			if !doubleQuote {
				if singleQuote {
					gotWord = true
				}
				singleQuote = !singleQuote
				continue
			}
		}
		gotWord = true
		buf.WriteRune(r)
	}
	if buf.Len() > 0 {
		out = append(out, buf.String())
	}
	return out
}
