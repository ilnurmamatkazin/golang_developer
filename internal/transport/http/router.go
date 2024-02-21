package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Выбрал легковесный и 100% совместимый со страндартным роутером.
// На текущем месте используем gin, но для тестового задания посчитал,
// что он будет излишним.
func (s *HttpServer) NewRouter() *chi.Mux {

	r := chi.NewRouter()

	// зададим встроенные middleware, чтобы улучшить стабильность приложения
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	initPrometheus(r)

	r.Route("/objects", func(r chi.Router) {
		r.Get("/{key}", s.getObject)
		r.With(middleware.AllowContentType("application/json")).Put("/{key}", s.createObject)
	})

	r.Route("/probes", func(r chi.Router) {
		r.Get("/liveness", s.liveness)
		r.Get("/readiness", s.readiness)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("route does not exist"))
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
		w.Write([]byte("method is not valid"))
	})

	return r

}
