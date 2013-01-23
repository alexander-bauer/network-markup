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
	stringsSEPARATOR = []string{",", "and", ", and"}
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

	// Loop through to tokenize
	tok := s.Scan()
	for tok != scanner.EOF {
		// Create a token, that we'll use
		t := &token{Literal: s.TokenText()}
		if matches(t.Literal, stringsNEWLINE) {
			t.Id = tokenNEWLINE
		} else if matches(t.Literal, stringsSEPARATOR) {
			t.Id = tokenSEPARATOR
		} else if matches(t.Literal, stringsCONNTO) {
			t.Id = tokenCONNTO
		} else {
			// Default to assuming that it's a reference
			t.Id = tokenNODE
		}
		println(t.Literal, t.Id)
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
