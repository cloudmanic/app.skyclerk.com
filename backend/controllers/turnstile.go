//
// Date: 2026-09-09
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2026 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// turnstileClient bounds verification latency and allows isolated HTTP mocks in tests.
var turnstileClient = &http.Client{Timeout: 10 * time.Second}

// verifyTurnstile validates single-use challenges before protected actions run.
// Reject missing configuration, outages, and tokens issued for other actions or sites.
func verifyTurnstile(request *http.Request, token, action string, hostnames ...string) bool {
	// Read the secret only on the server; local development also requires verification.
	secret := os.Getenv("TURNSTILE_SECRET_KEY")
	if secret == "" || strings.TrimSpace(token) == "" || len(token) > 2048 {
		return false
	}
	form := url.Values{"secret": {secret}, "response": {token}}
	// Propagate request cancellation and encode credentials in the POST body.
	verification, err := http.NewRequestWithContext(request.Context(), http.MethodPost,
		"https://challenges.cloudflare.com/turnstile/v0/siteverify", strings.NewReader(form.Encode()))
	if err != nil {
		return false
	}
	// Use Cloudflare's supported form encoding.
	verification.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// Bound network delays so an unavailable verification service cannot hang a form submission.
	response, err := turnstileClient.Do(verification)
	if err != nil {
		return false
	}
	// Release the connection regardless of verification outcome.
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return false
	}
	var result struct {
		Success  bool   `json:"success"`
		Action   string `json:"action"`
		Hostname string `json:"hostname"`
	}
	// Limit response size and reject malformed JSON or an unexpected widget action.
	if json.NewDecoder(io.LimitReader(response.Body, 16384)).Decode(&result) != nil || !result.Success || result.Action != action {
		return false
	}
	// Bind the challenge to the configured site instead of trusting the request Host.
	for _, hostname := range hostnames {
		if hostname != "" && result.Hostname == hostname {
			return true
		}
	}
	return false
}

/* End File */
