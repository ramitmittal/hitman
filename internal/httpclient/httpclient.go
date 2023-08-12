package httpclient

import (
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/ramitmittal/hitman/internal/parser"
)

type HitResult struct {
	Err             error
	RequestHeaders  []string
	ResponseHeaders []string
	ResponseBody    string
}

func formatRequest(req *http.Request) []string {
	reqHeaders := []string{
		req.Method + " " + req.URL.String(),
	}
	for h, v := range req.Header {
		for _, vv := range v {
			reqHeaders = append(reqHeaders, h+" : "+vv)
		}
	}
	return reqHeaders
}

func formatResponseHeaders(res *http.Response) []string {
	// equal length is a good starting point even though
	// headers slice may not have the same length as response header map
	// as response headers are flattened
	resHeaders := make([]string, 0, len(res.Header))
	for h, v := range res.Header {
		for _, vv := range v {
			resHeaders = append(resHeaders, h+" : "+vv)
		}
	}
	sort.Strings(resHeaders)
	return append([]string{res.Status}, resHeaders...)
}

var (
	flagInsecureSkipVerify = "insecure"
)

// Perform an HTTP request based on the command text
func Hit(text string) (hr *HitResult) {
	hr = &HitResult{}

	parserResult, err := parser.Parse([]byte(text))
	if err != nil {
		hr.Err = errors.New("please enter a valid query")
		return
	}

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	if _, prs := parserResult.Flags[flagInsecureSkipVerify]; prs {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	url := parserResult.Url
	if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
		url = "https://" + url
	}

	req, err := http.NewRequest(parserResult.Method, url, nil)
	if err != nil {
		hr.Err = err
		return
	}
	for k, v := range parserResult.Headers {
		req.Header.Add(k, v)
	}
	if req.Header.Get("Host") != "" {
		// Go httpClient treats Host header specially
		// Set it on the request directly
		req.Host = req.Header.Get("Host")
	}

	hr.RequestHeaders = formatRequest(req)

	res, err := client.Do(req)
	if err != nil {
		hr.Err = err
		return
	}
	hr.ResponseHeaders = formatResponseHeaders(res)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		hr.Err = err
		return
	}
	_ = res.Body.Close()

	// Scan the response body for non-printable characters
	// and set the response body to a dummy value to not mess up the terminal
	idx := 0
	for idx < len(body) {
		r, size := utf8.DecodeRune(body[idx:])
		if r != utf8.RuneError {
			idx += size
		} else if body[idx] == byte(10) {
			// Apparently line feeds are also not printable
			idx += 1
		} else {
			hr.ResponseBody = "\n\nRESPONSE CONTAINS NON-PRINTABLE CHARACTERS.\n"
			return
		}
	}

	hr.ResponseBody = string(body)
	return
}
