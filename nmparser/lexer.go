package nmparser

import (
	"io"
	"text/scanner"
)

const (
	tokenNEWLINE   = iota // Newline, semicolon, or other delimeter
	tokenSEPARATOR        // Comma, or "and" or ", and"
	tokenCONNTO           // "is connected to"
	tokenNODE             // Node reference such as "Alice"
)

var (
	stringsNEWLINE   = []string{"\n", ";", "."}
	stringsSEPARATOR = []string{",", "and"}
	stringsCONNTO    = []string{"is", "connected", "to"}
)

type token struct {
	Id      int    // Token type
	Literal string // Original string
}

func tokenize(src io.Reader) (tokens []*token) {
	// Initialize the Scanner
	var s scanner.Scanner
	s.Init(src)
	s.Mode = scanner.ScanIdents | scanner.ScanStrings

	// To support spaces in labels
	var terminated = true

	// Loop through to tokenize
	tok := s.Scan()
	for tok != scanner.EOF {
		// Create a token, that we'll use
		t := &token{Literal: s.TokenText()}
		if matches(t.Literal, stringsNEWLINE) {
			terminated = true
			t.Id = tokenNEWLINE
		} else if matches(t.Literal, stringsSEPARATOR) {
			terminated = true
			t.Id = tokenSEPARATOR
		} else if matches(t.Literal, stringsCONNTO) {
			terminated = true
			t.Id = tokenCONNTO
		} else {
			// Default to assuming that it's a reference
			if !terminated {
				// If it turns out that we're appending to the
				// previous token, then we just append to its label,
				// but then we have to make sure that we don't add
				// another token afterward. Several hours of headache
				// were spent discovering this.
				tokens[len(tokens)-1].Literal += " " + t.Literal
				tok = s.Scan()
				continue
			} else {
				t.Id = tokenNODE
			}
			terminated = false
		}
		tokens = append(tokens, t)

		tok = s.Scan()
	}
	return
}

func matches(str string, match []string) bool {
	for _, m := range match {
		if str == m {
			return true
		}
	}
	return false
}
