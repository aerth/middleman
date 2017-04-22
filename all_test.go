package middleman

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testhandler struct{}

func (t *testhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello, world"))
}

func TestAllMiddleware(t *testing.T) {
	var o = new(testhandler)
	h := CORS([]string{"http://example.com", "https://example.com"}, SingleHost("example.com", o))
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	req.Header.Set("Origin", "http://example.com")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	resp := w.Result()

	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	fmt.Println(resp.Header)
	fmt.Println(string(b))
	if string(b) != "hello, world" {
		t.Fail()
	}

	if resp.Header["Vary"][0] != "Origin" {
		t.Fail()
		t.Logf(`Wanted Header["Vary"] == "Origin", got %q`, resp.Header["Vary"][0])
	}
	// tbc
}
