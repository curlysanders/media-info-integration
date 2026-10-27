package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeAuthenticator struct {
	user string
	err  error
}

func (f fakeAuthenticator) AuthenticatedUser(context.Context, *http.Request) (string, error) {
	return f.user, f.err
}

func TestRuntimeHandlerReturnsAuthenticatedPackageProbe(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/mediainfo-poc/runtime", nil)
	request.Header.Set("X-Forwarded-User", "forged")
	request.Header.Set("X-Dsm-Server-Addr", "192.0.2.1")
	response := httptest.NewRecorder()
	runtimeHandler(fakeAuthenticator{user: "dsm-user"}, "abc123").ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if got := response.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("Cache-Control = %q, want no-store", got)
	}
	if !strings.Contains(response.Body.String(), `"user":"dsm-user"`) {
		t.Fatalf("response does not contain authenticated user: %s", response.Body.String())
	}
	if strings.Contains(response.Body.String(), "forged") {
		t.Fatalf("response trusts forged forwarded identity: %s", response.Body.String())
	}
}

func TestRuntimeHandlerFailsClosedWhenAuthenticationFails(t *testing.T) {
	response := httptest.NewRecorder()
	runtimeHandler(fakeAuthenticator{err: errors.New("helper failed")}, "abc123").ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/mediainfo-poc/runtime", nil))
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusUnauthorized)
	}
}

func TestRuntimeHandlerRejectsUnsupportedMethods(t *testing.T) {
	response := httptest.NewRecorder()
	runtimeHandler(fakeAuthenticator{user: "dsm-user"}, "abc123").ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/mediainfo-poc/runtime", nil))
	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusMethodNotAllowed)
	}
	if got := response.Header().Get("Allow"); got != http.MethodGet {
		t.Fatalf("Allow = %q, want GET", got)
	}
}
