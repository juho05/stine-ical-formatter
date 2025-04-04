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
	s.metrics.Visit()
	s.renderer.render(w, r, http.StatusOK, "main", templateData{})
}

func (s *Server) handlePostMainPage(w http.ResponseWriter, r *http.Request) {
	if err := tollbooth.LimitByRequest(s.limiter, w, r); err != nil {
		log.Warn("rate limit reached")
		s.metrics.FailureRateLimit()
		s.renderer.render(w, r, http.StatusTooManyRequests, "main", templateData{
			ErrorMessageKey: "error.rate-limit",
		})
		return
	}
	err := r.ParseMultipartForm(5e6) // 5 MB
	if err != nil {
		var maxBytesError *http.MaxBytesError
		errKey := "error.unexpected"
		if errors.As(err, &maxBytesError) {
			errKey = "files too large (sum must be <5MB)"
			s.metrics.FailureTooLarge()
			log.Warnf("uploaded file is too large: content length: %d", r.ContentLength)
		} else {
			s.metrics.FailureParseForm()
			log.Errorf("failed parse multipart form: %s", err)
		}
		s.renderer.render(w, r, http.StatusBadRequest, "main", templateData{
			ErrorMessageKey: errKey,
		})
		return
	}

	if len(r.MultipartForm.File["files"]) == 0 {
		log.Warn("no file uploaded")
		s.metrics.FailureNoFiles()
		s.renderer.render(w, r, http.StatusBadRequest, "main", templateData{
			ErrorMessageKey: "error.no-files",
		})
		return
	}
	files := make([]io.Reader, 0, len(r.MultipartForm.File["files"]))
	for _, v := range r.MultipartForm.File["files"] {
		if !slices.Contains([]string{".ical", ".ics", ".ifb", ".icalendar"}, strings.ToLower(filepath.Ext(v.Filename))) {
			log.Warnf("wrong extension: %s", v.Filename)
			s.metrics.FailureWrongFile()
			s.renderer.render(w, r, http.StatusBadRequest, "main", templateData{
				ErrorMessageKey: "error.not-ics",
			})
			return
		}
		f, err := v.Open()
		if err != nil {
			s.metrics.FailureOther()
			log.Errorf("failed to open uploaded file: %s", err)
			s.renderer.render(w, r, http.StatusInternalServerError, "main", templateData{
				ErrorMessageKey: "error.failed-to-open",
			})
			return
		}
		files = append(files, f)
		defer f.Close()
	}

	output, err := formatter.Format(files)
	if err != nil {
		s.metrics.FailureFormat()
		log.Errorf("failed to format files: %s", err)
		s.renderer.render(w, r, http.StatusInternalServerError, "main", templateData{
			ErrorMessageKey: "error.format-failed",
		})
		return
	}

	w.Header().Set("Content-Type", "text/calendar")
	w.Header().Set("Content-Disposition", `attachment; filename="calendar.ics"`)
	_, _ = w.Write(output)
	s.metrics.Success()
}
