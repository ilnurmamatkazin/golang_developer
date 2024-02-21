package http

import (
	"test/internal/model"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Инициализируем метрики для prometheus
func initPrometheus(r *chi.Mux) {

	prometheus.MustRegister(
		model.ObjectAddCounter,
		model.ObjectGetCounter,
		model.ObjectDelCounter,
		model.ObjectCountGauge,
	)

	r.Handle("/metrics", promhttp.Handler())

}
