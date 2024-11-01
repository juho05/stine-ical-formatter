package web

import (
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"slices"
	"strings"

	"github.com/didip/tollbooth/v7"
	"github.com/juho05/log"
	"github.com/juho05/stine-ical-formatter/formatter"
)

func (s *Server) handleGetMainPage(w http.ResponseWriter, r *http.Request) {
	s.renderer.render(w, r, http.StatusOK, "main", templateData{})
}

// 15 MB
const maxFileSize = 15e6

func (s *Server) handlePostMainPage(w http.ResponseWriter, r *http.Request) {
	if err := tollbooth.LimitByRequest(s.limiter, w, r); err != nil {
		s.renderer.render(w, r, http.StatusOK, "main", templateData{
			ErrorMessage: "rate limit exceeded",
		})
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)
	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		var maxBytesError *http.MaxBytesError
		errMsg := "an unexpected error occured"
		if errors.As(err, &maxBytesError) {
			errMsg = "files too large (sum must be <15MB)"
		} else {
			log.Errorf("failed parse multipart form: %w", err)
		}
		s.renderer.render(w, r, http.StatusOK, "main", templateData{
			ErrorMessage: errMsg,
		})
		return
	}

	if len(r.MultipartForm.File["files"]) == 0 {
		s.renderer.render(w, r, http.StatusOK, "main", templateData{
			ErrorMessage: "no files selected",
		})
		return
	}
	files := make([]io.Reader, 0, len(r.MultipartForm.File["files"]))
	for _, v := range r.MultipartForm.File["files"] {
		if !slices.Contains([]string{".ical", ".ics", ".ifb", ".icalendar"}, strings.ToLower(filepath.Ext(v.Filename))) {
			s.renderer.render(w, r, http.StatusOK, "main", templateData{
				ErrorMessage: "at least one file is not an iCalendar file (extension: .ics)",
			})
			return
		}
		f, err := v.Open()
		if err != nil {
			log.Errorf("failed to open uploaded file: %w", err)
			s.renderer.render(w, r, http.StatusOK, "main", templateData{
				ErrorMessage: "file(s) could not be opened",
			})
			return
		}
		files = append(files, f)
		defer f.Close()
	}

	output, err := formatter.Format(files)
	if err != nil {
		log.Errorf("failed to format files: %w", err)
		s.renderer.render(w, r, http.StatusOK, "main", templateData{
			ErrorMessage: err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "text/calendar")
	w.Header().Set("Content-Disposition", `attachment; filename="calendar.ics"`)
	_, _ = w.Write(output)
}
