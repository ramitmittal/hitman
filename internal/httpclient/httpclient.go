package httpclient

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ramitmittal/hitman/internal/parser"
)

type HitResult struct {
	Err error
	req *http.Request
	res *http.Response
}

func (hr *HitResult) RequestHeaders() []string {
	headers := []string{
		hr.req.Method + " " + hr.req.URL.String(),
	}
	for h, v := range hr.req.Header {
		for _, vv := range v {
			headers = append(headers, h+" : "+vv)
		}
	}
	return headers
}

func (hr *HitResult) ResponseHeaders() []string {
	headers := []string{
		hr.res.Status,
	}
	for h, v := range hr.res.Header {
		for _, vv := range v {
			headers = append(headers, h+" : "+vv)
		}
	}
	return headers
}

func (hr *HitResult) String() string {
	var sb strings.Builder
	for _, v := range hr.RequestHeaders() {
		sb.WriteString(v)
		sb.WriteRune('\n')
	}
	sb.WriteRune('\n')
	for _, v := range hr.ResponseHeaders() {
		sb.WriteString(v)
		sb.WriteRune('\n')
	}
	return sb.String()
}

// Perform an HTTP request based on the command text
func Hit(text string) HitResult {
	parserResult, err := parser.Parse([]byte(text))
	if err != nil {
		return HitResult{
			Err: errors.New("please enter a valid query"),
		}
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
		return HitResult{Err: err}
	}
	for k, v := range parserResult.GetHeaders() {
		req.Header.Add(k, v)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return HitResult{Err: err, req: req}
	}
	_ = res.Body.Close()

	return HitResult{req: req, res: res}
}
