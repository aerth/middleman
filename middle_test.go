package middleman

import (
	"bytes"
	"io"
	"log"
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
func TestWrapper(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<html><body>Hello World!</body></html>")
	}
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", 0)
	// add middleware
	handler = Hello(Log(logger, handler))

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	resp.Body.Close()
	if resp.Header.Get("Hello") != "World" {
		t.Fail()
		t.Logf(`Wanted "World", got: %q`, resp.Header.Get("Hello"))

	}

	logcat := buf.String()
	switch logcat {
	case "":
		t.Fail()
		t.Logf("Wanted log entries in buf, got empty string")
	default:
		t.Log("Got log entry:", logcat)
	}

}
