package main

import (
	"controllers"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIndex(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	controllers.Index(res, req)

	// Test if we get 200
	if res.Code != http.StatusOK {
		t.Fatalf("Response body did not contain expected %v:\n\tbody: %v",
			"200", res.Code)
	}

	body := res.Body.String()
	if !strings.Contains(body, "<!DOCTYPE html>") {
		t.Fatalf("Response body did not contain expected %v:\n\tbody: %v",
			"<!DOCTYPE html>", body)
	}
}
