package middleman

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCors(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<html><body>Hello World!</body></html>")
	}

	// add middleware
	handler := CORS([]string{"http://example.com", "https://example.com"}, http.HandlerFunc(h))

	// do request
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	req.Header.Add("Origin", "http://example.com")
	
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// response
	resp := w.Result()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Log(err.Error())
	}
	resp.Body.Close()

	println(string(b))
	fmt.Println(resp.Header)

	if acao := resp.Header["Access-Control-Allow-Origin"]; acao != nil && acao[0] != "http://example.com" {
		t.Fail()
		t.Logf(`Wanted "Access-Control-Allow-Origin" to be "http://example.com", got: %q`, acao)
	}

	if vary := resp.Header["Vary"]; vary != nil && vary[0] != "Origin" {
		t.Fail()
		t.Logf(`Wanted Header["Vary"] == "Origin", got: %q`, vary)
	}
	
}
