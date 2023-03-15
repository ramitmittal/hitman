package parser

import (
	"errors"
	"strings"
	"unicode/utf8"
)

//go:generate go run golang.org/x/tools/cmd/goyacc -l -o parser.go parser.y

type Result struct {
	method  string
	url     string
	headers map[string]string
	flags   map[string]string
}

func (r Result) GetMethod() string {
	return r.method
}
func (r Result) GetURL() string {
	return r.url
}
func (r Result) GetHeaders() map[string]string {
	return r.headers
}

func (r Result) GetFlags() map[string]string {
	return r.flags
}

func Parse(input []byte) (Result, error) {
	l := &lex{
		input: input,
	}
	_ = yyParse(l)
	return l.result, l.err
}

type lex struct {
	input  []byte
	result Result
	err    error

	position int
}

func (l *lex) Lex(lval *yySymType) int {
	r, size := utf8.DecodeRune(l.input[l.position:])
	l.position += size

	if size == 0 {
		return 0
	}
	if r == ' ' || r == '\n' {
		// discard spaces and newlines
		return l.Lex(lval)
	}
	if r == ':' {
		return int(r)
	}
	if r == '#' {
		// discard everything till \n
		for {
			r1, size1 := utf8.DecodeRune(l.input[l.position:])
			l.position += size1
			if size1 == 0 {
				return 0
			}
			if r1 == '\n' {
				break
			}
		}
		return l.Lex(lval)
	}
	if r == '"' {
		// gather everything till closing "
		var str strings.Builder
		for {
			r1, size1 := utf8.DecodeRune(l.input[l.position:])
			l.position += size1
			if size1 == 0 {
				return 0
			}
			if r1 == '"' {
				lval.val = str.String()
				return S
			}
			str.WriteRune(r1)
		}
	}
	if r == '-' {
		// gather everything till space or newline
		var str strings.Builder
		for {
			r1, size1 := utf8.DecodeRune(l.input[l.position:])
			l.position += size1
			if size1 == 0 {
				return 0
			}
			if r1 == ' ' || r1 == '\n' {
				lval.val = str.String()
				return Flag
			}
			str.WriteRune(r1)
		}
	}

	// gather everything till one of the aforementioned characters
	var str strings.Builder
	str.WriteRune(r)

	for {
		r1, size1 := utf8.DecodeRune(l.input[l.position:])
		if r1 == '\n' || r1 == ' ' || r1 == ':' || size1 == 0 {
			lval.val = str.String()
			return S
		}

		l.position += size1
		str.WriteRune(r1)
	}
}

func (l *lex) Error(s string) {
	l.err = errors.New(s)
}

func merge(x, y map[string]string) map[string]string {
	if y == nil {
		return x
	}
	for k, v := range x {
		y[k] = v
	}
	return y
}

func mapOf(k, v string) map[string]string {
	return map[string]string{
		k: v,
	}
}
