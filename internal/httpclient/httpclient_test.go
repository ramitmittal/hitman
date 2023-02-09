package httpclient

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customHeader := r.Header["X-Custom-Header"]
		if len(customHeader) != 1 {
			w.WriteHeader(http.StatusBadRequest)
		} else if customHeader[0] != "Hello, World!" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	input := fmt.Sprintf(`GET "%s" X-Custom-Header: "Hello, World!"`, server.URL)
	hr := Hit(input)

	if hr.Err != nil {
		t.Fail()
	} else if hr.res.StatusCode != http.StatusOK {
		t.Fail()
	}
}
