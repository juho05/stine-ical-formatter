package web

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
)

type renderer struct {
	templates map[string]*template.Template
}

type templateData struct {
	ErrorMessage string
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

func (r *renderer) render(w http.ResponseWriter, status int, page string, data templateData) {
	t, ok := r.templates[page]
	if !ok {
		serverError(w, fmt.Errorf("template %s does not exist", page))
		return
	}

	buf := &bytes.Buffer{}
	err := t.ExecuteTemplate(buf, "base", data)
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
