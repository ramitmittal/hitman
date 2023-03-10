package parser

import "testing"

func TestValidQuotes(t *testing.T) {
	input := `GET "https://www.ramitmittal.com"
Accept-Encoding: "gzip, br"`
	inputBytes := []byte(input)

	if v, err := Parse(inputBytes); err != nil {
		t.Fail()
	} else if v.url != "https://www.ramitmittal.com" {
		t.Fail()
	} else if v.headers["Accept-Encoding"] != "gzip, br" {
		t.Fail()
	}
}

func TestInvalidQuotes(t *testing.T) {
	input := `GET https://www.ramitmittal.com
Accept-Encoding: gzip, "br"`
	inputBytes := []byte(input)

	if _, err := Parse(inputBytes); err == nil {
		t.Fail()
	}
}

func TestValidInputs(t *testing.T) {
	var tests = []struct {
		name  string
		input string
	}{
		{
			"simple request with one header",
			`GET www.ramitmittal.com XXX: hello`,
		},
		{
			"first line can be a comment",
			`#POST www.ramitmittal.com
GET www.ramitmittal.com
XXX: hello`,
		},
		{
			"header line can be a comment",
			`GET www.ramitmittal.com
# XXX: world
XXX: hello`,
		},
		{
			"comment line can be empty",
			`GET www.ramitmittal.com
#
XXX: hello`,
		},
		{
			"extra newlines are allowed",
			`GET www.ramitmittal.com

XXX: hello

`,
		},
		{
			"multiple headers are allowed",
			`GET www.ramitmittal.com
Cache-Control: no-cache
XXX: hello`,
		},
		{
			"headers may have 0 spaces around :",
			`GET www.ramitmittal.com
XXX:hello`,
		},
		{
			"headers may have multiple spaces around :",
			`GET www.ramitmittal.com
XXX  :  hello`,
		},
		{
			"trailing spaces in first line",
			`GET www.ramitmittal.com   
XXX:hello`,
		},
		{
			"inline comments",
			`GET www.ramitmittal.com
Cache-Control: no-cache # a comment
XXX: hello # another comment`,
		},
		{
			"multiple headers on same line",
			`GET www.ramitmittal.com
Cache-Control: no-cache XXX:hello`,
		},
		{
			"multiple # for comments",
			`GET www.ramitmittal.com
## XXX: world
XXX: hello`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inputBytes := []byte(test.input)

			if v, err := Parse(inputBytes); err != nil {
				t.Fail()
			} else if v.method != "GET" {
				t.Fail()
			} else if v.headers["XXX"] != "hello" {
				t.Fail()
			}
		})
	}
}

func TestInvalidInputs(t *testing.T) {
	var tests = []struct {
		name  string
		input string
	}{
		{"Multiple methods", `GET POST www.ramitmittal.com`},
		{"Request scheme in URL without quotes", `GET https://www.ramitmittal.com`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inputBytes := []byte(test.input)

			if _, err := Parse(inputBytes); err == nil {
				t.Fail()
			}
		})
	}
}
