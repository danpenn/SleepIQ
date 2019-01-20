package sleepiq

import (
	"net/http"
	"testing"
)

func TestHttpGetSuccess(t *testing.T) {
	var cookies []*http.Cookie

	responseBytes, err := httpGet("http://msn.com", cookies, getHeaders())
	if err != nil {
		t.Error("request failed - expected success", err)
		return
	}

	if len(responseBytes) <= 0 {
		t.Error("request failed - no data was returned")
	}
}

func TestHttpGetBadUrl(t *testing.T) {
	var cookies []*http.Cookie

	_, err := httpGet("foo://invalid.com", cookies, getHeaders())
	if err == nil {
		t.Error("request succeeded - expected failure", err)
		return
	}
}

func TestHttpPutBadUrl(t *testing.T) {
	testBytes := []byte("testing")
	var cookies []*http.Cookie
	_, _, err := httpPut("foo://invalid.com", testBytes, cookies)
	if err == nil {
		t.Error("request succeeded - expected failure", err)
		return
	}
}
