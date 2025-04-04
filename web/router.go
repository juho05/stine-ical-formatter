package web

import (
	"context"
	"embed"
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
		w.Header().Set("Permissions-Policy", "geolocation=(), camera=(), microphone=()")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-site")
		w.Header().Set("Permissions-Policy", "interest-cohort=()")
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
	r.Get("/static/css/tailwind.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, staticFS, "static/css/tailwind.css")
	})
}
