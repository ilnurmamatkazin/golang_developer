package http

import (
	"net/http"

	"github.com/go-chi/render"
)

// Метод сообщает, что наш сервис запущен
func (s *HttpServer) liveness(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)

}

// Метод сообщает, что наш сервис готов к работе
func (s *HttpServer) readiness(w http.ResponseWriter, r *http.Request) {

	if s.stor.Stor.IsNotInit() {
		render.Status(r, http.StatusServiceUnavailable)
		render.PlainText(w, r, http.StatusText(http.StatusServiceUnavailable))

		return
	}

	w.WriteHeader(http.StatusOK)

}
