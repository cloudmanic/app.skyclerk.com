//
// Date: 2026-09-09
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2026 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
)

// TestContactUs verifies that only valid support challenges can send notifications,
// and that malformed requests, outages, and failed delivery never report success.
func TestContactUs(t *testing.T) {
	cases := []struct {
		name, token, verification  string
		verificationStatus, status int
	}{
		{"missing token", "", "", 0, 400},
		{"blank token", "   ", "", 0, 400},
		{"oversized token", strings.Repeat("a", 2049), "", 0, 400},
		{"missing secret", "test-token", "", 0, 400},
		{"invalid token", "test-token", `{"success":false}`, 200, 400},
		{"expired or replayed", "test-token", `{"success":false,"error-codes":["timeout-or-duplicate"]}`, 200, 400},
		{"signup token", "test-token", `{"success":true,"action":"register","hostname":"skyclerk.com"}`, 200, 400},
		{"missing action", "test-token", `{"success":true,"hostname":"skyclerk.com"}`, 200, 400},
		{"wrong hostname", "test-token", `{"success":true,"action":"support_contact","hostname":"other.example"}`, 200, 400},
		{"app hostname", "test-token", `{"success":true,"action":"support_contact","hostname":"app.skyclerk.com"}`, 200, 400},
		{"missing hostname", "test-token", `{"success":true,"action":"support_contact"}`, 200, 400},
		{"malformed verification", "test-token", `{`, 200, 400},
		{"service unavailable", "test-token", `{"success":true,"action":"support_contact","hostname":"skyclerk.com"}`, 503, 400},
		{"network failure", "test-token", "", -1, 400},
		{"malformed request", "test-token", "", 0, 400},
		{"invalid email", "test-token", "", 0, 400},
		{"missing message", "test-token", "", 0, 400},
		{"oversized request", "test-token", "", 0, 400},
		{"verified apex", "test-token", `{"success":true,"action":"support_contact","hostname":"skyclerk.com"}`, 200, 204},
		{"verified www", "test-token", `{"success":true,"action":"support_contact","hostname":"www.skyclerk.com"}`, 200, 204},
		{"email delivery failure", "test-token", `{"success":true,"action":"support_contact","hostname":"skyclerk.com"}`, 200, 503},
	}
	for _, tc := range cases {
		// Exercise the actual HTTP handler while isolating every external service.
		t.Run(tc.name, func(t *testing.T) {
			// Require verification even in development, without reading real credentials.
			t.Setenv("APP_ENV", "local")
			t.Setenv("TURNSTILE_SECRET_KEY", "test-secret")
			if tc.name == "missing secret" {
				// Missing server configuration must not provide a bypass.
				t.Setenv("TURNSTILE_SECRET_KEY", "")
			}
			// Verify the outgoing Cloudflare request using the shared in-process transport.
			verificationCalls := mockTurnstileResponse(t, tc.verificationStatus, tc.verification)
			previousSend, previousNotify := supportSendEmail, supportNotify
			// Restore delivery functions when this case finishes, including assertion failures.
			t.Cleanup(func() {
				supportSendEmail, supportNotify = previousSend, previousNotify
			})
			emailCalls := 0
			notifications := make(chan string, 1)
			// Capture email content without contacting a mail provider.
			supportSendEmail = func(to, replyTo, subject, body string, attachments []string) error {
				emailCalls++
				// Confirm the existing recipient and reply behavior, with user HTML escaped.
				st.Expect(t, to, "help@skyclerk.com")
				st.Expect(t, replyTo, "visitor@example.com")
				st.Expect(t, subject, "[Website Contact]: New request from Test Visitor")
				st.Expect(t, strings.Contains(body, "&lt;script&gt;example&lt;/script&gt;"), true)
				st.Expect(t, strings.Contains(body, "test-token"), false)
				if tc.name == "email delivery failure" {
					// Simulate a delivery outage after Cloudflare has consumed the token.
					return errors.New("mail provider unavailable")
				}
				return nil
			}
			// Capture asynchronous notifications without sending any Slack messages.
			supportNotify = func(channel, message string) (string, error) {
				notifications <- channel + ":" + message
				return "ok", nil
			}
			post := map[string]string{
				"fullName": "Test Visitor", "email": "visitor@example.com", "phone": "555-0100",
				"message": "<script>example</script>", "turnstile_token": tc.token,
			}
			if tc.name == "invalid email" {
				post["email"] = "invalid"
			}
			if tc.name == "missing message" {
				post["message"] = ""
			}
			if tc.name == "oversized request" {
				// Reject oversized input before Cloudflare or notification work begins.
				post["message"] = strings.Repeat("a", 65536)
			}
			// Encode the same JSON contract used by the marketing site's Vue form.
			body, err := json.Marshal(post)
			if err != nil {
				// Stop immediately if the test request cannot be constructed.
				t.Fatal(err)
			}
			if tc.name == "malformed request" {
				body = []byte("{")
			}
			controller := &Controller{}
			// Route through Gin so JSON binding and response codes match production.
			router := gin.New()
			// Mount the public support endpoint without any authentication middleware.
			router.POST("/support/contact-us", controller.ContactUs)
			// Capture the HTTP response locally instead of starting a server.
			writer := httptest.NewRecorder()
			// Send the request through the actual public handler.
			router.ServeHTTP(writer, httptest.NewRequest(http.MethodPost, "/support/contact-us", bytes.NewReader(body)))
			// Every failure must stop before notifications, except a failed email attempt itself.
			st.Expect(t, writer.Code, tc.status)
			st.Expect(t, *verificationCalls > 0, tc.verificationStatus != 0)
			st.Expect(t, emailCalls > 0, tc.status == 204 || tc.status == 503)
			if tc.status == 204 {
				// Wait for the asynchronous notification, with a short failure deadline.
				select {
				case notice := <-notifications:
					// Only verified and delivered support requests should reach Slack.
					st.Expect(t, notice, "#events:Skyclerk Website Support Request: Email: visitor@example.com")
				case <-time.After(time.Second):
					// Fail if successful delivery did not schedule the existing notification.
					t.Fatal("expected a support notification")
				}
			} else {
				select {
				case notice := <-notifications:
					// Rejected requests and email failures cannot notify Slack.
					t.Fatalf("unexpected support notification: %s", notice)
				default:
				}
			}
			if tc.status == 400 && tc.name != "malformed request" && tc.name != "invalid email" && tc.name != "missing message" && tc.name != "oversized request" {
				// Security errors must be recognizable so the frontend can offer a fresh challenge.
				st.Expect(t, strings.Contains(writer.Body.String(), `"code":"turnstile_failed"`), true)
			}
		})
	}
}

/* End File */
