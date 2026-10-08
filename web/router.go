package web

import (
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed static
var staticFS embed.FS

func generateNonce(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "nonce", generateToken(32))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nonce := r.Context().Value("nonce").(string)
		w.Header().Set("Content-Security-Policy", fmt.Sprintf("default-src 'none'; script-src 'self' 'nonce-%s'; img-src 'self'; style-src 'self' 'nonce-%s'; font-src 'self';", nonce, nonce))
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Permissions-Policy", "geolocation=(), camera=(), microphone=(), interest-cohort=()")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-site")
		next.ServeHTTP(w, r)
	})
}

func registerMiddlewares(r chi.Router) {
	r.Use(generateNonce)
	r.Use(securityHeaders)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestSize(6e6)) // 6 MB
	r.Use(middleware.CleanPath)
}

func (s *Server) registerRoutes(r chi.Router) {
	r.Get("/", s.handleGetMainPage)
	r.Post("/", s.handlePostMainPage)
	r.Post("/export", s.handlePostExport)
	r.Get("/metrics", s.metrics.ServeHTTP)
	r.Get("/static/css/tailwind.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", `"`+cssVersion+`"`)
		if r.URL.Query().Get("v") == cssVersion {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			w.Header().Set("Cache-Control", "no-cache")
		}
		http.ServeFileFS(w, r, staticFS, "static/css/tailwind.css")
	})
}

var cssVersion = mustHashStaticFile("static/css/tailwind.css")

func mustHashStaticFile(name string) string {
	data, err := staticFS.ReadFile(name)
	if err != nil {
		panic(fmt.Errorf("read %s: %w", name, err))
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:8])
}
