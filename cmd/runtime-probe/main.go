package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultAuthenticator = "/usr/syno/synoman/webman/modules/authenticate.cgi"
	maxIdentityBytes     = 256
	maxEnvironmentValue  = 16 * 1024
	maxRequestBodyBytes  = 1024
	authTimeout          = 2 * time.Second
	authWaitDelay        = 250 * time.Millisecond
)

var authenticationEnvironmentNames = map[string]struct{}{
	"CONTENT_LENGTH": {}, "CONTENT_TYPE": {}, "GATEWAY_INTERFACE": {}, "HTTPS": {}, "HTTP_COOKIE": {}, "HTTP_HOST": {}, "PATH_INFO": {},
	"QUERY_STRING": {}, "REMOTE_ADDR": {}, "REQUEST_METHOD": {}, "REQUEST_URI": {}, "SCRIPT_NAME": {},
	"SERVER_ADDR": {}, "SERVER_NAME": {}, "SERVER_PORT": {}, "SERVER_PROTOCOL": {}, "HTTP_X_SYNO_TOKEN": {},
}

var buildRevision = "dev"

type authenticator interface {
	AuthenticatedUser(context.Context, []string, []byte) (string, error)
}

type cgiAuthenticator struct{ path string }

func (a cgiAuthenticator) AuthenticatedUser(ctx context.Context, environment []string, body []byte) (string, error) {
	command := authenticationCommand(ctx, a.path, environment, body)
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

func authenticationCommand(ctx context.Context, path string, environment []string, body []byte) *exec.Cmd {
	command := exec.CommandContext(ctx, path)
	command.Env = authenticationEnvironment(environment)
	command.Stdin = bytes.NewReader(body)
	command.WaitDelay = authWaitDelay
	return command
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

type routingPolicy struct {
	AllowedUsers   []string `json:"allowedUsers"`
	AllowedOrigins []string `json:"allowedOrigins"`
}

func (p routingPolicy) permitsUser(username string) bool { return contains(p.AllowedUsers, username) }
func (p routingPolicy) permitsOrigin(origin string) bool { return contains(p.AllowedOrigins, origin) }

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

func loadRoutingPolicy(path string) (routingPolicy, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return routingPolicy{}, err
	}
	var policy routingPolicy
	decoder := json.NewDecoder(strings.NewReader(string(contents)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&policy); err != nil {
		return routingPolicy{}, err
	}
	if err := requireEOF(decoder); err != nil {
		return routingPolicy{}, err
	}
	if len(policy.AllowedUsers) == 0 || len(policy.AllowedOrigins) == 0 {
		return routingPolicy{}, errors.New("routing policy has an empty allowlist")
	}
	return policy, nil
}

type probeRequest struct {
	Operation string `json:"operation"`
}

type response struct {
	Runtime  string `json:"runtime"`
	Revision string `json:"revision"`
	UID      int    `json:"uid"`
	GID      int    `json:"gid"`
	User     string `json:"user"`
}

func probeResponse(ctx context.Context, auth authenticator, policy routingPolicy, environment []string, body io.Reader, revision string, uid, gid int) (int, []byte) {
	if cgiValue(environment, "REQUEST_METHOD") != "POST" {
		return 405, []byte(`{"error":"method not allowed"}`)
	}
	if cgiValue(environment, "CONTENT_TYPE") != "application/json" {
		return 400, []byte(`{"error":"invalid request"}`)
	}
	if cgiValue(environment, "HTTP_X_SYNO_TOKEN") == "" {
		return 401, []byte(`{"error":"authentication required"}`)
	}
	length, present, err := requestLength(environment)
	if err != nil {
		return 400, []byte(`{"error":"invalid request"}`)
	}
	if !present {
		return 400, []byte(`{"error":"invalid request"}`)
	}
	if length > maxRequestBodyBytes {
		return 413, []byte(`{"error":"request too large"}`)
	}
	requestBody, err := readRequestBody(body, length)
	if err != nil {
		return 400, []byte(`{"error":"invalid request"}`)
	}

	timedContext, cancel := context.WithTimeout(ctx, authTimeout)
	defer cancel()
	username, err := auth.AuthenticatedUser(timedContext, environment, requestBody)
	if err != nil {
		return 401, []byte(`{"error":"authentication required"}`)
	}
	if !policy.permitsUser(username) || !policy.permitsOrigin(cgiValue(environment, "HTTP_ORIGIN")) {
		return 403, []byte(`{"error":"request forbidden"}`)
	}
	if err := validateProbeRequest(requestBody); err != nil {
		return 400, []byte(`{"error":"invalid request"}`)
	}

	encoded, err := json.Marshal(response{Runtime: "go-cgi", Revision: revision, UID: uid, GID: gid, User: username})
	if err != nil {
		return 500, []byte(`{"error":"response encoding failed"}`)
	}
	return 200, append(encoded, '\n')
}

func readRequestBody(body io.Reader, length int) ([]byte, error) {
	contents := make([]byte, length)
	if _, err := io.ReadFull(body, contents); err != nil {
		return nil, errors.New("invalid request body")
	}
	return contents, nil
}

func validateProbeRequest(body []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()
	var request probeRequest
	if err := decoder.Decode(&request); err != nil {
		return err
	}
	if err := requireEOF(decoder); err != nil || request.Operation != "route-probe" {
		return errors.New("invalid route probe request")
	}
	return nil
}

func requireEOF(decoder *json.Decoder) error {
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("unexpected trailing JSON")
		}
		return err
	}
	return nil
}

func requestLength(environment []string) (int, bool, error) {
	value := cgiValue(environment, "CONTENT_LENGTH")
	if value == "" {
		return 0, false, nil
	}
	var length int
	if _, err := fmt.Sscanf(value, "%d", &length); err != nil || length < 0 || fmt.Sprintf("%d", length) != value {
		return 0, true, errors.New("invalid content length")
	}
	return length, true, nil
}

func cgiValue(environment []string, name string) string {
	prefix := name + "="
	for _, item := range environment {
		if value, found := strings.CutPrefix(item, prefix); found {
			return value
		}
	}
	return ""
}

func defaultPolicyPath() string {
	executable, err := os.Executable()
	if err != nil {
		return "routing-policy.json"
	}
	return filepath.Join(filepath.Dir(executable), "routing-policy.json")
}

func main() {
	revision := flag.String("revision", buildRevision, "prototype package revision")
	authenticatorPath := flag.String("authenticator", defaultAuthenticator, "DSM authentication helper")
	policyPath := flag.String("policy", defaultPolicyPath(), "routing policy JSON path")
	flag.Parse()
	policy, err := loadRoutingPolicy(*policyPath)
	if err != nil {
		writeResponse(500, []byte(`{"error":"routing unavailable"}`))
		return
	}
	status, body := probeResponse(context.Background(), cgiAuthenticator{path: *authenticatorPath}, policy, os.Environ(), os.Stdin, *revision, os.Getuid(), os.Getgid())
	writeResponse(status, body)
}

func writeResponse(status int, body []byte) {
	fmt.Printf("Status: %d %s\r\n", status, statusText(status))
	fmt.Print("Content-Type: application/json; charset=utf-8\r\n")
	fmt.Print("Cache-Control: no-store\r\n\r\n")
	_, _ = os.Stdout.Write(body)
}

func statusText(status int) string {
	switch status {
	case 200:
		return "OK"
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 403:
		return "Forbidden"
	case 405:
		return "Method Not Allowed"
	case 413:
		return "Payload Too Large"
	default:
		return "Internal Server Error"
	}
}
