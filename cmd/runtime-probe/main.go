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
	AuthenticatedUser(context.Context, http.Header, string, string) (string, error)
}

type cgiAuthenticator struct{ path string }

func (a cgiAuthenticator) AuthenticatedUser(ctx context.Context, header http.Header, remoteAddr, serverAddr string) (string, error) {
	command := exec.CommandContext(ctx, a.path)
	command.Env = []string{
		"PATH=/usr/bin:/bin",
		"GATEWAY_INTERFACE=CGI/1.1",
		"REQUEST_METHOD=GET",
		"SCRIPT_NAME=/mediainfo-poc/runtime",
		"HTTP_COOKIE=" + header.Get("Cookie"),
		"REMOTE_ADDR=" + remoteAddr,
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
		username, err := auth.AuthenticatedUser(ctx, request.Header, request.RemoteAddr, request.Header.Get("X-Dsm-Server-Addr"))
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
