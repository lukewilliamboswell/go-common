package common

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORS(t *testing.T) {

	var (
		err error
		ts  *httptest.Server
		res *http.Response
	)

	ts = httptest.NewServer(CORS(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "Hello, World!")
	})))
	defer ts.Close()

	res, err = http.Get(ts.URL)
	if err != nil {
		t.Log("should return a response")
		t.Fail()
	}

	if res.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Log("should have 'Access-Control-Allow-Origin' header set")
		t.Fail()
	}
}
