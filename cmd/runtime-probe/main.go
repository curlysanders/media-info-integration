package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

const (
	defaultSocket        = "/var/packages/mediainfo-poc/var/runtime.sock"
	defaultAuthenticator = "/usr/syno/synoman/webman/modules/authenticate.cgi"
	maxIdentityBytes     = 256
	authTimeout          = 2 * time.Second
)

type authenticator interface {
	AuthenticatedUser(context.Context, *http.Request) (string, error)
}

type cgiAuthenticator struct{ path string }

func (a cgiAuthenticator) AuthenticatedUser(ctx context.Context, request *http.Request) (string, error) {
	serverAddr := request.Header.Get("X-Dsm-Server-Addr")
	command := exec.CommandContext(ctx, a.path)
	command.Env = []string{
		"PATH=/usr/bin:/bin",
		"GATEWAY_INTERFACE=CGI/1.1",
		"REQUEST_METHOD=" + request.Method,
		"SCRIPT_NAME=" + request.URL.Path,
		"REQUEST_URI=" + request.RequestURI,
		"QUERY_STRING=" + request.URL.RawQuery,
		"SERVER_PROTOCOL=" + request.Proto,
		"SERVER_NAME=" + request.Host,
		"SERVER_PORT=5000",
		"HTTP_HOST=" + request.Host,
		"HTTP_COOKIE=" + request.Header.Get("Cookie"),
		"REMOTE_ADDR=" + request.RemoteAddr,
		"SERVER_ADDR=" + serverAddr,
		"HTTP_X_DSM_SERVER_ADDR=" + serverAddr,
	}
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

type response struct {
	Runtime  string `json:"runtime"`
	Revision string `json:"revision"`
	UID      int    `json:"uid"`
	GID      int    `json:"gid"`
	User     string `json:"user"`
}

func runtimeHandler(auth authenticator, revision string) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Cache-Control", "no-store")
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		if request.Method != http.MethodGet {
			writer.Header().Set("Allow", http.MethodGet)
			http.Error(writer, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
		ctx, cancel := context.WithTimeout(request.Context(), authTimeout)
		defer cancel()
		username, err := auth.AuthenticatedUser(ctx, request)
		if err != nil {
			http.Error(writer, `{"error":"authentication required"}`, http.StatusUnauthorized)
			return
		}
		writer.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(writer).Encode(response{Runtime: "go-fastcgi", Revision: revision, UID: os.Getuid(), GID: os.Getgid(), User: username})
	})
}

func main() {
	socket := flag.String("socket", defaultSocket, "Unix socket for the FastCGI server")
	revision := flag.String("revision", "dev", "prototype package revision")
	authenticatorPath := flag.String("authenticator", defaultAuthenticator, "DSM authentication helper")
	flag.Parse()
	socketDirectory := filepath.Dir(*socket)
	info, err := os.Stat(socketDirectory)
	if err != nil {
		log.Fatalf("inspect socket directory: %v", err)
	}
	if !info.IsDir() {
		log.Fatalf("socket directory is not a directory: %s", socketDirectory)
	}
	if err := os.Remove(*socket); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("remove stale socket: %v", err)
	}
	listener, err := net.Listen("unix", *socket)
	if err != nil {
		log.Fatalf("listen on socket: %v", err)
	}
	if err := os.Chmod(*socket, 0666); err != nil {
		_ = listener.Close()
		log.Fatalf("set socket mode: %v", err)
	}
	defer func() { _ = listener.Close(); _ = os.Remove(*socket) }()
	server := http.NewServeMux()
	server.Handle("/mediainfo-poc/runtime", runtimeHandler(cgiAuthenticator{path: *authenticatorPath}, *revision))
	terminated := make(chan os.Signal, 1)
	signal.Notify(terminated, syscall.SIGINT, syscall.SIGTERM)
	go func() { <-terminated; _ = listener.Close() }()
	if err := fcgi.Serve(listener, server); err != nil && !errors.Is(err, net.ErrClosed) {
		fmt.Fprintf(os.Stderr, "serve FastCGI: %v\n", err)
		os.Exit(1)
	}
}
