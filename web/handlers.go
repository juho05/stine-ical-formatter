package web

import "net/http"

func (s *Server) handleGetMainPage(w http.ResponseWriter, r *http.Request) {
	s.renderer.render(w, http.StatusOK, "main", templateData{})
}
