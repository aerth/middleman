package middleman

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type oh struct{}
func (o *oh) ServeHTTP(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("hello"))
} 

func TestAllMiddleware(t *testing.T) {
	var o = new(oh)
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
	if string(b) != "hello" {
		t.Fail()
	}

	// tbc

}
