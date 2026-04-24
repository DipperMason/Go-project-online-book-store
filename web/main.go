package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type app struct {
	authBase    *url.URL
	profileBase *url.URL
	client      *http.Client
}

func main() {
	authURL, err := url.Parse(getEnv("AUTH_URL", "http://localhost:8001"))
	if err != nil {
		log.Fatalf("invalid AUTH_URL: %v", err)
	}

	profileURL, err := url.Parse(getEnv("PROFILE_URL", "http://localhost:8003"))
	if err != nil {
		log.Fatalf("invalid PROFILE_URL: %v", err)
	}

	a := &app{
		authBase:    authURL,
		profileBase: profileURL,
		client:      &http.Client{},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/auth/register", a.proxyAuthRegister)
	mux.HandleFunc("/api/auth/login", a.proxyAuthLogin)
	mux.HandleFunc("/api/profile", a.proxyProfile)
	mux.HandleFunc("/api/profile/books", a.proxyProfileBooks)
	mux.HandleFunc("/api/profile/activity", a.proxyProfileActivity)
	mux.HandleFunc("/", serveStatic)

	addr := getEnv("WEB_ADDR", ":8088")
	log.Printf("web ui started at %s", addr)
	log.Printf("auth proxy -> %s", authURL.String())
	log.Printf("profile proxy -> %s", profileURL.String())
	if err = http.ListenAndServe(addr, logging(mux)); err != nil {
		log.Fatalf("web server error: %v", err)
	}
}

func (a *app) proxyAuthRegister(w http.ResponseWriter, r *http.Request) {
	a.proxy(w, r, a.authBase, "/api/v1/register")
}

func (a *app) proxyAuthLogin(w http.ResponseWriter, r *http.Request) {
	a.proxy(w, r, a.authBase, "/api/v1/auth")
}

func (a *app) proxyProfile(w http.ResponseWriter, r *http.Request) {
	a.proxy(w, r, a.profileBase, "/api/v1/profile")
}

func (a *app) proxyProfileBooks(w http.ResponseWriter, r *http.Request) {
	a.proxy(w, r, a.profileBase, "/api/v1/profile/books")
}

func (a *app) proxyProfileActivity(w http.ResponseWriter, r *http.Request) {
	a.proxy(w, r, a.profileBase, "/api/v1/profile/activity")
}

func (a *app) proxy(w http.ResponseWriter, r *http.Request, base *url.URL, path string) {
	target := *base
	target.Path = path
	target.RawQuery = r.URL.RawQuery

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request", http.StatusBadRequest)
		return
	}

	req, err := http.NewRequestWithContext(r.Context(), r.Method, target.String(), strings.NewReader(string(body)))
	if err != nil {
		http.Error(w, "failed to create proxy request", http.StatusInternalServerError)
		return
	}

	copyHeaders(r.Header, req.Header)
	req.Header.Del("Origin")

	resp, err := a.client.Do(req)
	if err != nil {
		http.Error(w, "upstream unavailable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	copyHeaders(resp.Header, w.Header())
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

func copyHeaders(src http.Header, dst http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

func serveStatic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		path = "index.html"
	}

	if strings.Contains(path, "..") {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(".", path)
	if _, err := os.Stat(fullPath); err != nil {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, fullPath)
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
