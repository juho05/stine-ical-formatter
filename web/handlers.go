package web

import "net/http"

func (s *Server) handleGetMainPage(w http.ResponseWriter, r *http.Request) {
	s.renderer.render(w, r, http.StatusOK, "main", templateData{})
}

func (s *Server) handlePostMainPage(w http.ResponseWriter, r *http.Request) {
	s.renderer.render(w, r, http.StatusOK, "main", templateData{})
}
