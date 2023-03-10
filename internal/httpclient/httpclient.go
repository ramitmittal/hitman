package httpclient

import (
	"errors"
	"io"
	"net/http"
	"strings"

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
	resHeaders := []string{
		res.Status,
	}
	for h, v := range res.Header {
		for _, vv := range v {
			resHeaders = append(resHeaders, h+" : "+vv)
		}
	}
	return resHeaders
}

func (hr *HitResult) String() string {
	var sb strings.Builder
	for _, v := range hr.RequestHeaders {
		sb.WriteString(v)
		sb.WriteRune('\n')
	}
	sb.WriteRune('\n')
	for _, v := range hr.ResponseHeaders {
		sb.WriteString(v)
		sb.WriteRune('\n')
	}
	sb.WriteRune('\n')
	sb.WriteString(hr.ResponseBody)
	return sb.String()
}

// Perform an HTTP request based on the command text
func Hit(text string) (hr HitResult) {
	parserResult, err := parser.Parse([]byte(text))
	if err != nil {
		hr.Err = errors.New("please enter a valid query")
		return
	}

	http.DefaultClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	url := parserResult.GetURL()
	if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
		url = "https://" + url
	}

	req, err := http.NewRequest(parserResult.GetMethod(), url, nil)
	if err != nil {
		hr.Err = err
		return
	}
	for k, v := range parserResult.GetHeaders() {
		req.Header.Add(k, v)
	}
	if req.Header.Get("Host") != "" {
		// Go httpClient treats Host header specially
		// Set it on the request directly and remove it from headers
		req.Host = req.Header.Get("Host")
		req.Header.Del("Host")
	}

	hr.RequestHeaders = formatRequest(req)

	res, err := http.DefaultClient.Do(req)
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
	hr.ResponseBody = string(body)

	return
}
