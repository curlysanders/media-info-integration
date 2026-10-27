package main

import (
	"context"
	"errors"
	"strings"
	"testing"
)

type fakeAuthenticator struct {
	user string
	err  error
}

func (f fakeAuthenticator) AuthenticatedUser(context.Context, []string) (string, error) {
	return f.user, f.err
}

func TestProbeResponseReturnsDSMAuthenticatedPackageProbe(t *testing.T) {
	status, body := probeResponse(context.Background(), fakeAuthenticator{user: "dsm-user"}, []string{
		"REQUEST_METHOD=GET",
		"HTTP_X_FORWARDED_USER=forged",
	}, "abc123", 123, 456)
	if status != 200 {
		t.Fatalf("status = %d, want 200", status)
	}
	if !strings.Contains(string(body), `"user":"dsm-user"`) {
		t.Fatalf("response does not contain authenticated user: %s", body)
	}
	if strings.Contains(string(body), "forged") {
		t.Fatalf("response trusts forged forwarded identity: %s", body)
	}
}

func TestProbeResponseFailsClosedWhenAuthenticationFails(t *testing.T) {
	status, _ := probeResponse(context.Background(), fakeAuthenticator{err: errors.New("helper failed")}, []string{"REQUEST_METHOD=GET"}, "abc123", 123, 456)
	if status != 401 {
		t.Fatalf("status = %d, want 401", status)
	}
}

func TestProbeResponseRejectsUnsupportedMethods(t *testing.T) {
	status, body := probeResponse(context.Background(), fakeAuthenticator{user: "dsm-user"}, []string{"REQUEST_METHOD=POST"}, "abc123", 123, 456)
	if status != 405 {
		t.Fatalf("status = %d, want 405", status)
	}
	if string(body) != `{"error":"method not allowed"}` {
		t.Fatalf("body = %s", body)
	}
}

func TestAuthenticationEnvironmentKeepsOnlyDocumentedCGIContext(t *testing.T) {
	filtered := strings.Join(authenticationEnvironment([]string{
		"REQUEST_METHOD=GET", "HTTP_COOKIE=session-value", "REMOTE_ADDR=192.0.2.1", "SERVER_ADDR=192.0.2.2",
		"HTTP_X_SYNO_TOKEN=dsm-csrf-token", "HTTP_X_FORWARDED_USER=forged", "AWS_SECRET_ACCESS_KEY=not-for-helper",
	}), "\n")
	for _, expected := range []string{"PATH=/usr/bin:/bin", "REQUEST_METHOD=GET", "HTTP_COOKIE=session-value", "REMOTE_ADDR=192.0.2.1", "SERVER_ADDR=192.0.2.2", "HTTP_X_SYNO_TOKEN=dsm-csrf-token"} {
		if !strings.Contains(filtered, expected) {
			t.Fatalf("environment missing %q: %s", expected, filtered)
		}
	}
	for _, forbidden := range []string{"HTTP_X_FORWARDED_USER", "AWS_SECRET_ACCESS_KEY"} {
		if strings.Contains(filtered, forbidden) {
			t.Fatalf("environment contains %q: %s", forbidden, filtered)
		}
	}
}

func TestAuthenticationEnvironmentDropsOversizedValues(t *testing.T) {
	filtered := strings.Join(authenticationEnvironment([]string{"HTTP_COOKIE=" + strings.Repeat("x", maxEnvironmentValue+1)}), "\n")
	if strings.Contains(filtered, "HTTP_COOKIE=") {
		t.Fatalf("environment retains oversized cookie")
	}
}
