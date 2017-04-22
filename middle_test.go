package middleman

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHello(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<html><body>Hello World!</body></html>")
	}

	// add middleware
	handler = Hello(handler)

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	resp.Body.Close()
	if resp.Header.Get("Hello") != "World" {
		t.Fail()
		t.Logf(`Wanted "World", got: %q`, resp.Header.Get("Hello"))

	}
}
