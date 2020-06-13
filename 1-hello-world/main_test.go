package main

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestHandleRoot(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080", nil)
	w := httptest.NewRecorder()
	HandleRoot(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if respStatusCode := resp.StatusCode; respStatusCode != 200 {
		t.Errorf("resp.StatusCode == %d, want %d", respStatusCode, 200)
	}

	if respBody := string(body); respBody != "Hello world!" {
		t.Errorf("resp.Body == %s, want %s", respBody, "Hello world!")
	}
}
