package main

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"testing"
)

type fakeAuthenticator struct {
	user string
	err  error
}

func (f fakeAuthenticator) AuthenticatedUser(context.Context, []string, []byte) (string, error) {
	return f.user, f.err
}

type recordingAuthenticator struct {
	environment []string
	body        []byte
}

func (a *recordingAuthenticator) AuthenticatedUser(_ context.Context, environment []string, body []byte) (string, error) {
	a.environment = append([]string(nil), environment...)
	a.body = append([]byte(nil), body...)
	return "test1", nil
}

var testPolicy = routingPolicy{
	AllowedUsers:   []string{"test1"},
	AllowedOrigins: []string{"https://nas.example.test:5001", "https://mediainfo.example.test"},
}

func validEnvironment() []string {
	return []string{
		"REQUEST_METHOD=POST", "CONTENT_TYPE=application/json", "CONTENT_LENGTH=27",
		"HTTP_X_SYNO_TOKEN=dsm-csrf-token", "HTTP_ORIGIN=https://nas.example.test:5001",
	}
}

func environmentForBody(body string) []string {
	environment := validEnvironment()
	environment[2] = "CONTENT_LENGTH=" + strconv.Itoa(len(body))
	return environment
}

func TestProbeResponseReturnsDSMAuthenticatedPackageProbe(t *testing.T) {
	environment := append(validEnvironment(), "HTTP_X_FORWARDED_USER=forged")
	status, body := probeResponse(context.Background(), fakeAuthenticator{user: "test1"}, testPolicy, environment, strings.NewReader(`{"operation":"route-probe"}`), "abc123", 123, 456)
	if status != 200 {
		t.Fatalf("status = %d, want 200", status)
	}
	if !strings.Contains(string(body), `"user":"test1"`) {
		t.Fatalf("response does not contain authenticated user: %s", body)
	}
	if strings.Contains(string(body), "forged") {
		t.Fatalf("response trusts forged forwarded identity: %s", body)
	}
}

func TestProbeResponsePassesBoundedCGIContextAndBodyToAuthenticator(t *testing.T) {
	auth := &recordingAuthenticator{}
	environment := validEnvironment()
	body := `{"operation":"route-probe"}`
	status, _ := probeResponse(context.Background(), auth, testPolicy, environment, strings.NewReader(body), "abc123", 123, 456)
	if status != 200 {
		t.Fatalf("status = %d, want 200", status)
	}
	if got := string(auth.body); got != body {
		t.Fatalf("authentication body = %q, want %q", got, body)
	}
	if got := cgiValue(auth.environment, "CONTENT_LENGTH"); got != "27" {
		t.Fatalf("authentication CONTENT_LENGTH = %q, want 27", got)
	}
}

func TestProbeResponseFailsClosedWhenAuthenticationFails(t *testing.T) {
	status, _ := probeResponse(context.Background(), fakeAuthenticator{err: errors.New("helper failed")}, testPolicy, validEnvironment(), strings.NewReader(`{"operation":"route-probe"}`), "abc123", 123, 456)
	if status != 401 {
		t.Fatalf("status = %d, want 401", status)
	}
}

func TestProbeResponseRejectsMissingDSMToken(t *testing.T) {
	environment := validEnvironment()
	environment[3] = "HTTP_X_SYNO_TOKEN="
	status, _ := probeResponse(context.Background(), fakeAuthenticator{user: "test1"}, testPolicy, environment, strings.NewReader(`{"operation":"route-probe"}`), "abc123", 123, 456)
	if status != 401 {
		t.Fatalf("status = %d, want 401", status)
	}
}

func TestProbeResponseRejectsDeniedUserAndOrigin(t *testing.T) {
	status, _ := probeResponse(context.Background(), fakeAuthenticator{user: "test2"}, testPolicy, validEnvironment(), strings.NewReader(`{"operation":"route-probe"}`), "abc123", 123, 456)
	if status != 403 {
		t.Fatalf("denied user status = %d, want 403", status)
	}
	environment := validEnvironment()
	environment[4] = "HTTP_ORIGIN=https://evil.example.test"
	status, _ = probeResponse(context.Background(), fakeAuthenticator{user: "test1"}, testPolicy, environment, strings.NewReader(`{"operation":"route-probe"}`), "abc123", 123, 456)
	if status != 403 {
		t.Fatalf("denied origin status = %d, want 403", status)
	}
}

func TestProbeResponseRejectsUnsupportedMethodsAndInvalidBodies(t *testing.T) {
	environment := validEnvironment()
	environment[0] = "REQUEST_METHOD=GET"
	status, _ := probeResponse(context.Background(), fakeAuthenticator{user: "test1"}, testPolicy, environment, strings.NewReader(`{"operation":"route-probe"}`), "abc123", 123, 456)
	if status != 405 {
		t.Fatalf("method status = %d, want 405", status)
	}
	for _, body := range []string{"", `{"operation":"wrong"}`, `{"operation":"route-probe","extra":true}`, `{"operation":"route-probe"} {}`} {
		status, _ = probeResponse(context.Background(), fakeAuthenticator{user: "test1"}, testPolicy, environmentForBody(body), strings.NewReader(body), "abc123", 123, 456)
		if status != 400 {
			t.Fatalf("body %q status = %d, want 400", body, status)
		}
	}
}

func TestProbeResponseRejectsOversizedBodyAndLength(t *testing.T) {
	overstated := validEnvironment()
	overstated[2] = "CONTENT_LENGTH=1025"
	status, _ := probeResponse(context.Background(), fakeAuthenticator{user: "test1"}, testPolicy, overstated, strings.NewReader(`{"operation":"route-probe"}`), "abc123", 123, 456)
	if status != 413 {
		t.Fatalf("declared length status = %d, want 413", status)
	}
	oversizedBody := strings.Repeat("x", maxRequestBodyBytes+1)
	status, _ = probeResponse(context.Background(), fakeAuthenticator{user: "test1"}, testPolicy, environmentForBody(oversizedBody), strings.NewReader(oversizedBody), "abc123", 123, 456)
	if status != 413 {
		t.Fatalf("body status = %d, want 413", status)
	}
	malformed := validEnvironment()
	malformed[2] = "CONTENT_LENGTH=not-a-number"
	status, _ = probeResponse(context.Background(), fakeAuthenticator{user: "test1"}, testPolicy, malformed, strings.NewReader(`{"operation":"route-probe"}`), "abc123", 123, 456)
	if status != 400 {
		t.Fatalf("malformed length status = %d, want 400", status)
	}
	missing := validEnvironment()
	missing[2] = "CONTENT_LENGTH="
	status, _ = probeResponse(context.Background(), fakeAuthenticator{user: "test1"}, testPolicy, missing, strings.NewReader(`{"operation":"route-probe"}`), "abc123", 123, 456)
	if status != 400 {
		t.Fatalf("missing length status = %d, want 400", status)
	}
	truncated := validEnvironment()
	truncated[2] = "CONTENT_LENGTH=28"
	status, _ = probeResponse(context.Background(), fakeAuthenticator{user: "test1"}, testPolicy, truncated, strings.NewReader(`{"operation":"route-probe"}`), "abc123", 123, 456)
	if status != 400 {
		t.Fatalf("truncated body status = %d, want 400", status)
	}
}

func TestAuthenticationEnvironmentKeepsOnlyDocumentedCGIContext(t *testing.T) {
	filtered := strings.Join(authenticationEnvironment([]string{
		"REQUEST_METHOD=POST", "CONTENT_TYPE=application/json", "CONTENT_LENGTH=27", "HTTP_COOKIE=session-value", "REMOTE_ADDR=192.0.2.1", "SERVER_ADDR=192.0.2.2",
		"HTTP_X_SYNO_TOKEN=dsm-csrf-token", "HTTP_X_FORWARDED_USER=forged", "AWS_SECRET_ACCESS_KEY=not-for-helper",
	}), "\n")
	for _, expected := range []string{"PATH=/usr/bin:/bin", "REQUEST_METHOD=POST", "CONTENT_TYPE=application/json", "CONTENT_LENGTH=27", "HTTP_COOKIE=session-value", "REMOTE_ADDR=192.0.2.1", "SERVER_ADDR=192.0.2.2", "HTTP_X_SYNO_TOKEN=dsm-csrf-token"} {
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

func TestAuthenticationCommandBoundsWaitForLeakedOutputPipes(t *testing.T) {
	command := authenticationCommand(context.Background(), "/bin/true", []string{"HTTP_COOKIE=session-value"}, []byte(`{"operation":"route-probe"}`))
	if command.WaitDelay != authWaitDelay {
		t.Fatalf("WaitDelay = %s, want %s", command.WaitDelay, authWaitDelay)
	}
	if got := strings.Join(command.Env, "\n"); !strings.Contains(got, "HTTP_COOKIE=session-value") {
		t.Fatalf("authentication command lost bounded CGI environment: %s", got)
	}
}
