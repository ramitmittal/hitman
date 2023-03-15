package parser

import "testing"

func TestValidQuotes(t *testing.T) {
	input := `GET "https://www.ramitmittal.com"
Accept-Encoding: "gzip, br"`
	inputBytes := []byte(input)

	if v, err := Parse(inputBytes); err != nil {
		t.Fail()
	} else if v.Url != "https://www.ramitmittal.com" {
		t.Fail()
	} else if v.Headers["Accept-Encoding"] != "gzip, br" {
		t.Fail()
	}
}

func TestValidFlags(t *testing.T) {
	var tests = []struct {
		name  string
		input string
	}{
		{"One flag and one header", `GET www.ramitmittal.com Cache-Control: "no-cache" -flag1`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inputBytes := []byte(test.input)

			if v, err := Parse(inputBytes); err != nil {
				t.Fail()
			} else if len(v.Flags) != 1 {
				t.Fail()
			} else if len(v.Headers) != 1 {
				t.Fail()
			}
		})
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
		{
			"flags after headers",
			`GET www.ramitmittal.com XXX: hello -flag1`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inputBytes := []byte(test.input)

			if v, err := Parse(inputBytes); err != nil {
				t.Fail()
			} else if v.Method != "GET" {
				t.Fail()
			} else if v.Headers["XXX"] != "hello" {
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
		{"Flags before headers", `GET www.ramitmittal.com -flag1 Cache-Control: "no-cache"`},
		{"URL with : must be quoted", `GET https://www.ramitmittal.com`},
		{"Quotes inside header values are not supported", `GET www.ramitmittal.com Accept-Encoding: gzip, "br"`},
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
