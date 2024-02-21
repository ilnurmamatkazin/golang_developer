package http

import (
	"net/http"
	"test/internal/model"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// Метод возвращает объект клиенту, если объект не найден (или был удален из-за
// просрочки), то возвращаем 404
func (s *HttpServer) getObject(w http.ResponseWriter, r *http.Request) {

	key := chi.URLParam(r, "key")

	object, err := s.stor.Stor.Get(key)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
		return
	}

	model.ObjectGetCounter.Inc()

	if object != nil {
		render.Status(r, http.StatusOK)
		render.JSON(w, r, object)
	} else {
		render.Status(r, http.StatusNotFound)
		render.PlainText(w, r, http.StatusText(http.StatusNotFound))
	}

}

// Метод сохранения объекта в системе
func (s *HttpServer) createObject(w http.ResponseWriter, r *http.Request) {

	var (
		expires time.Time
		body    interface{}
		err     error
	)

	if r.Header.Get("Expires") != "" {
		if expires, err = time.Parse("2006-01-02 15:04:05", r.Header.Get("Expires")); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.PlainText(w, r, "В заголовке Expires указан неверный формат даты.")
			return
		}
	}

	if err = render.DecodeJSON(r.Body, &body); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, err.Error())
		return
	}

	key := chi.URLParam(r, "key")

	if err = s.stor.Stor.Add(key, body, expires); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
		return
	}

	model.ObjectAddCounter.Inc()

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, body)

}
