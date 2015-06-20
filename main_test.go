package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleIndexReturnsWithStatusOk(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	index(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("Response body did not contain expected %v:\n\tbody: %v",
			"200", res.Code)
	}
}
