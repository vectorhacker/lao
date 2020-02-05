package lao

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// Kind enum
type Kind int

// Kind enum
const (
	_ Kind = iota
	KindIdentifier
	KindKeyword
	KindInteger
	KindReal
	KindPeriod
	KindLogicalOperator
	KindArithmeticOperator
	KindRelationalOperator
	KindString
	KindAssignment
	KindEnd
)

// Token from tokenizer.
type Token struct {
	Kind   Kind
	Value  string
	Line   int
	Column int
}

// Tokenizer takes in a stream of input and produces tokens
type Tokenizer interface {
	Current() Token
	Next() bool
}

// NewTokenizer creates a new tokenizier
func NewTokenizer(r io.Reader) Tokenizer {
	buf := new(bytes.Buffer)

	buf.ReadFrom(r)

	return &tokenizer{
		buf:    buf,
		line:   1,
		column: 1,
	}
}

type tokenizer struct {
	buf      *bytes.Buffer
	ct       Token
	position int
	line     int
	column   int
}

func (t *tokenizer) Current() Token {
	return t.ct
}
func (t *tokenizer) Next() bool {
	if t.position >= t.buf.Len() {
		t.ct = Token{Kind: KindEnd}
		return false
	}

	t.skipWhitespaceAndNewLines()

	ch := t.buf.Bytes()[t.position]

	if unicode.IsLetter(rune(ch)) {
		t.recognizeKeywordsAndIdentifier()
	}

	if unicode.IsDigit(rune(ch)) {
		t.recognizeNumber()
	}

	if ch == '-' || ch == '+' {
		t.recognizeNumber()
	}

	if ch == '.' {
		t.recognizeOperatorsAndPeriods()
	}

	if ch == '=' {
		t.recognizeAssignment()
	}

	if ch == '"' {
		t.recognizeString()
	}

	return true
}

func (t *tokenizer) recognizeString() {
	pos := t.position + 1
	column := t.column

	stringFinsihed := false
	for pos < t.buf.Len() {
		pos++
		ch := t.buf.Bytes()[pos]
		if ch == '"' {
			stringFinsihed = true
			break
		}
	}

	if stringFinsihed {

		b := make([]byte, pos-t.position+1)

		read := copy(b, t.buf.Bytes()[t.position:])
		_ = read

		s := fmt.Sprintf("%s", b)

		t.ct = Token{
			Kind:   KindString,
			Value:  s,
			Column: column,
			Line:   t.line,
		}
		t.position += len(s)
		t.column += len(s)
	}
}

func isKeyword(s string) bool {

	switch s {
	case "print", "rem", "if", "read", "then", "end":
		return true
	}

	return false
}

func (t *tokenizer) recognizeKeywordsAndIdentifier() {

	pos := t.position
	column := t.column

	identifier := ""
	for pos < t.buf.Len() {
		ch := t.buf.Bytes()[pos]

		if !unicode.IsLetter(rune(ch)) {
			break
		}

		identifier += string(ch)
		pos++
	}

	if isKeyword(strings.ToLower(identifier)) {

		t.ct = Token{
			Kind:   KindKeyword,
			Line:   t.line,
			Column: column,
			Value:  identifier,
		}
	} else {

		t.ct = Token{
			Kind:   KindIdentifier,
			Line:   t.line,
			Column: column,
			Value:  identifier,
		}
	}

	t.position = pos
	t.column += len(identifier)

}

func (t *tokenizer) recognizeOperatorsAndPeriods() {

	value := ""
	line := t.line
	colunm := t.column
	// read till whitespace or newline
	for i := t.position; i < t.buf.Len(); i++ {
		ch := t.buf.Bytes()[i]

		if unicode.IsSpace(rune(ch)) {
			break
		}

		value += string(ch)
	}

	switch strings.ToLower(value) {
	case ".add.", ".sub.", ".mul.", ".div.":
		t.ct = Token{
			Kind:   KindArithmeticOperator,
			Value:  value,
			Column: colunm,
			Line:   line,
		}
		t.column += len(value)
		t.position += len(value)
	case ".gt.", ".lt.", ".ge.", ".le.", ".eq.", ".ne.":
		t.ct = Token{
			Kind:   KindRelationalOperator,
			Value:  value,
			Column: colunm,
			Line:   line,
		}
		t.column += len(value)
		t.position += len(value)
	case ".not.", ".and.", ".or.":
		t.ct = Token{
			Kind:   KindLogicalOperator,
			Value:  value,
			Column: colunm,
			Line:   line,
		}
		t.column += len(value)
		t.position += len(value)
	case ".":
		t.ct = Token{
			Kind:   KindPeriod,
			Value:  value,
			Column: colunm,
			Line:   line,
		}
		t.column += len(value)
		t.position += len(value)
	}
}

type numberState int

const (
	initial numberState = iota
	integer
	beginSignedNumber
	signedNumber
	beginNumberWithFractionalPart
	numberWithFractionalPart
	beginNumberWithExponent
	beginNumberWithSignedExponent
	numberWithExponent
	noNextState
)

func (t *tokenizer) recognizeNumber() {

	line := t.line
	column := t.column

	nextState := func(currentState numberState, ch byte) numberState {

		switch currentState {
		case initial:
			if unicode.IsDigit(rune(ch)) {
				return integer
			}
			if ch == '+' || ch == '-' {
				return beginSignedNumber
			}
		case signedNumber:
			if unicode.IsDigit(rune(ch)) {
				return integer
			}
		case beginSignedNumber:
			if unicode.IsDigit(rune(ch)) {
				return signedNumber
			}
		case integer:
			if unicode.IsDigit(rune(ch)) {
				return integer
			}

			if ch == '.' {
				return beginNumberWithFractionalPart
			}

			if unicode.ToLower(rune(ch)) == 'e' {
				return beginNumberWithExponent
			}
		case beginNumberWithFractionalPart:
			if unicode.IsDigit(rune(ch)) {
				return numberWithFractionalPart
			}
		case numberWithFractionalPart:
			if unicode.IsDigit(rune(ch)) {
				return numberWithFractionalPart
			}
			if unicode.ToLower(rune(ch)) == 'e' {
				return beginNumberWithExponent
			}
		case numberWithExponent:
			if unicode.IsDigit(rune(ch)) {
				return numberWithExponent
			}
			if ch == '.' {
				return numberWithExponent
			}
		case beginNumberWithExponent:
			if ch == '+' || ch == '-' {
				return beginNumberWithSignedExponent
			}

			if unicode.IsDigit(rune(ch)) {
				return numberWithExponent
			}

		case beginNumberWithSignedExponent:
			if unicode.IsDigit(rune(ch)) {
				return numberWithExponent
			}
		}

		return noNextState
	}

	acceptingStates := []numberState{numberWithExponent, integer, numberWithFractionalPart}

	has := func(states []numberState, state numberState) bool {
		for _, accepted := range states {
			if state == accepted {
				return true
			}
		}

		return false
	}

	run := func() (bool, numberState, string) {
		current := initial
		number := ""

		for i := t.position; i < t.buf.Len(); i++ {
			ch := t.buf.Bytes()[i]
			next := nextState(current, ch)

			if next == noNextState {
				break
			}
			number += string(ch)

			current = next
		}

		return has(acceptingStates, current), current, number
	}

	if isNumber, state, number := run(); isNumber {

		var kind Kind

		switch state {
		case integer:
			kind = KindInteger
		case numberWithExponent, numberWithFractionalPart:
			kind = KindReal
		}

		t.ct = Token{
			Kind:   kind,
			Value:  number,
			Line:   line,
			Column: column,
		}
		t.position += len(number)
		t.column += len(number)
	}

}

func (t *tokenizer) recognizeAssignment() {
	if t.buf.Bytes()[t.position] == '=' {
		t.ct = Token{
			Kind:   KindAssignment,
			Value:  "=",
			Line:   t.line,
			Column: t.column,
		}
		t.column++
		t.position++
	}
}

func (t *tokenizer) skipWhitespaceAndNewLines() {

	for t.position < t.buf.Len() &&
		(unicode.IsSpace(rune(t.buf.Bytes()[t.position])) ||
			t.buf.Bytes()[t.position] == '\n') {
		if t.buf.Bytes()[t.position] == '\n' {
			t.line++
			t.column = 1
		} else {
			t.column++
		}
		t.position++
	}
}
