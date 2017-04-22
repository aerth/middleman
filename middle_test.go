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

func TestBoolFunc(t *testing.T) {
	hf := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<html><body>Hello World!</body></html>")
	}

	// add middleware
	failall := func(w http.ResponseWriter, r *http.Request) bool {
		http.Error(w, "error", http.StatusMethodNotAllowed)
		return false
	}
	handler := IfThen(failall, http.HandlerFunc(hf))

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fail()
		t.Logf("wanted %v, got %v", resp.StatusCode, http.StatusMethodNotAllowed)
	}
	t.Logf("wanted %v, got %v", resp.StatusCode, http.StatusMethodNotAllowed)

}
func TestSingleHost(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<html><body>Hello World!</body></html>")
	}

	// add middleware
	handler := SingleHost("example.org", http.HandlerFunc(h))

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	resp.Body.Close()
	if resp.StatusCode != 403 {
		t.Fail()
		t.Log("Wanted 403, got", resp.StatusCode)
	}

	req = httptest.NewRequest("GET", "http://example.org/foo", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp = w.Result()
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fail()
		t.Log("Wanted 200, got", resp.StatusCode)
	}

}
