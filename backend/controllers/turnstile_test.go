//
// Date: 2026-09-09
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2026 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/nbio/st"
)

// turnstileTestTransport adapts a function into an isolated HTTP transport.
type turnstileTestTransport func(*http.Request) (*http.Response, error)

// RoundTrip invokes the fake verifier without changing the global HTTP transport.
func (transport turnstileTestTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return transport(request)
}

// mockTurnstileResponse validates outgoing form fields and supplies a deterministic
// response or network error. The previous client is restored after each test.
func mockTurnstileResponse(t *testing.T, status int, body string) *int {
	previous := turnstileClient
	calls := 0
	// Restore the production client even when an assertion fails.
	t.Cleanup(func() { turnstileClient = previous })
	// Keep all verification traffic in-process and assert the endpoint and payload.
	turnstileClient = &http.Client{Transport: turnstileTestTransport(func(request *http.Request) (*http.Response, error) {
		calls++
		st.Expect(t, request.URL.String(), "https://challenges.cloudflare.com/turnstile/v0/siteverify")
		st.Expect(t, request.Method, http.MethodPost)
		st.Expect(t, request.Header.Get("Content-Type"), "application/x-www-form-urlencoded")
		// Decode the POST fields to ensure neither secret nor token is omitted.
		if err := request.ParseForm(); err != nil {
			t.Fatal(err)
		}
		st.Expect(t, request.PostForm.Get("secret"), "test-secret")
		st.Expect(t, request.PostForm.Get("response"), "test-token")
		if status == -1 {
			return nil, errors.New("verification unavailable")
		}
		if status == 0 {
			t.Fatal("unexpected verification request")
		}
		return &http.Response{StatusCode: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
	})}
	return &calls
}

/* End File */
