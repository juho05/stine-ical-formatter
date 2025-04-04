package web

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type renderer struct {
	templates map[string]*template.Template
}

type templateData struct {
	ErrorMessageKey string
}

//go:embed templates/*.tmpl.html
var templateFS embed.FS

func newRenderer() (*renderer, error) {
	r := &renderer{
		templates: make(map[string]*template.Template),
	}
	err := r.loadTemplates(templateFS)
	if err != nil {
		return nil, fmt.Errorf("load templates: %w", err)
	}
	return r, nil
}

func (r *renderer) render(w http.ResponseWriter, req *http.Request, status int, page string, data templateData) {
	t, ok := r.templates[page]
	if !ok {
		serverError(w, fmt.Errorf("template %s does not exist", page))
		return
	}

	lang := langFromAcceptLanguageHeader(req.Header.Get("Accept-Language"))

	buf := &bytes.Buffer{}
	type tmplData struct {
		Nonce string
		templateData
		Translate func(key string) template.HTML
	}
	err := t.ExecuteTemplate(buf, "base", tmplData{
		Nonce:        req.Context().Value("nonce").(string),
		templateData: data,
		Translate: func(key string) template.HTML {
			return template.HTML(Translate(lang, key))
		},
	})
	if err != nil {
		serverError(w, err)
		return
	}

	w.WriteHeader(status)
	_, _ = buf.WriteTo(w)
}

func (r *renderer) loadTemplates(htmlFS fs.FS) error {
	pages, err := fs.Glob(htmlFS, "templates/*.tmpl.html")
	if err != nil {
		return fmt.Errorf("find html pages: %w", err)
	}

	for _, page := range pages {
		name := strings.TrimSuffix(filepath.Base(page), ".tmpl.html")

		t, err := template.New(name).ParseFS(htmlFS, "templates/base.tmpl.html")
		if err != nil {
			return fmt.Errorf("parse base.tmpl.html: %w", err)
		}

		t, err = t.ParseFS(htmlFS, page)
		if err != nil {
			return fmt.Errorf("parse %s: %w", page, err)
		}

		r.templates[name] = t
	}

	return nil
}

func langFromAcceptLanguageHeader(headerValue string) string {
	lang := "en"
	quality := float64(0)

	strs := strings.Split(headerValue, ",")
	for _, s := range strs {
		parts := strings.Split(s, ";")
		q := float64(1)
		if len(parts) > 1 {
			qStr := parts[1]
			qStr = strings.ReplaceAll(qStr, "q", "")
			qStr = strings.ReplaceAll(qStr, "=", "")
			qStr = strings.TrimSpace(qStr)
			_, err := strconv.ParseFloat(qStr, 64)
			if err == nil {
				q, _ = strconv.ParseFloat(qStr, 64)
			}
		}

		if q > quality {
			l := strings.TrimSpace(strings.Split(parts[0], "-")[0])
			if _, ok := translations[l]; ok {
				lang = l
				quality = q
			}
		}
	}

	return lang
}
