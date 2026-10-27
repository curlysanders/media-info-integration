package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	defaultAuthenticator = "/usr/syno/synoman/webman/modules/authenticate.cgi"
	maxIdentityBytes     = 256
	maxEnvironmentValue  = 16 * 1024
	authTimeout          = 2 * time.Second
)

var authenticationEnvironmentNames = map[string]struct{}{
	"GATEWAY_INTERFACE": {}, "HTTPS": {}, "HTTP_COOKIE": {}, "HTTP_HOST": {}, "PATH_INFO": {},
	"QUERY_STRING": {}, "REMOTE_ADDR": {}, "REQUEST_METHOD": {}, "REQUEST_URI": {}, "SCRIPT_NAME": {},
	"SERVER_ADDR": {}, "SERVER_NAME": {}, "SERVER_PORT": {}, "SERVER_PROTOCOL": {},
}

type authenticator interface {
	AuthenticatedUser(context.Context, []string) (string, error)
}

type cgiAuthenticator struct{ path string }

func (a cgiAuthenticator) AuthenticatedUser(ctx context.Context, environment []string) (string, error) {
	command := exec.CommandContext(ctx, a.path)
	command.Env = authenticationEnvironment(environment)
	output, err := command.Output()
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return "", errors.New("authentication helper timed out")
	}
	if err != nil {
		return "", errors.New("authentication helper failed")
	}
	if len(output) > maxIdentityBytes {
		return "", errors.New("authentication helper returned too much data")
	}
	username := strings.TrimSpace(string(output))
	if username == "" || strings.ContainsAny(username, "\r\n\x00") {
		return "", errors.New("request is not authenticated")
	}
	return username, nil
}

func authenticationEnvironment(environment []string) []string {
	filtered := []string{"PATH=/usr/bin:/bin"}
	for _, item := range environment {
		name, value, found := strings.Cut(item, "=")
		if !found || len(value) > maxEnvironmentValue {
			continue
		}
		if _, allowed := authenticationEnvironmentNames[name]; allowed {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

type response struct {
	Runtime  string `json:"runtime"`
	Revision string `json:"revision"`
	UID      int    `json:"uid"`
	GID      int    `json:"gid"`
	User     string `json:"user"`
}

func probeResponse(ctx context.Context, auth authenticator, environment []string, revision string, uid, gid int) (int, []byte) {
	if requestMethod(environment) != "GET" {
		return 405, []byte(`{"error":"method not allowed"}`)
	}
	timedContext, cancel := context.WithTimeout(ctx, authTimeout)
	defer cancel()
	username, err := auth.AuthenticatedUser(timedContext, environment)
	if err != nil {
		return 401, []byte(`{"error":"authentication required"}`)
	}
	body, err := json.Marshal(response{Runtime: "go-cgi", Revision: revision, UID: uid, GID: gid, User: username})
	if err != nil {
		return 500, []byte(`{"error":"response encoding failed"}`)
	}
	return 200, append(body, '\n')
}

func requestMethod(environment []string) string {
	for _, item := range environment {
		if method, found := strings.CutPrefix(item, "REQUEST_METHOD="); found {
			return method
		}
	}
	return ""
}

func main() {
	revision := flag.String("revision", "dev", "prototype package revision")
	authenticatorPath := flag.String("authenticator", defaultAuthenticator, "DSM authentication helper")
	flag.Parse()
	status, body := probeResponse(context.Background(), cgiAuthenticator{path: *authenticatorPath}, os.Environ(), *revision, os.Getuid(), os.Getgid())
	fmt.Printf("Status: %d %s\r\n", status, statusText(status))
	fmt.Print("Content-Type: application/json; charset=utf-8\r\n")
	fmt.Print("Cache-Control: no-store\r\n\r\n")
	_, _ = os.Stdout.Write(body)
}

func statusText(status int) string {
	switch status {
	case 200:
		return "OK"
	case 401:
		return "Unauthorized"
	case 405:
		return "Method Not Allowed"
	default:
		return "Internal Server Error"
	}
}
