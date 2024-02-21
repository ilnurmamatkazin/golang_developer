package model

import (
	"github.com/prometheus/client_golang/prometheus"
)

var ObjectAddCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "object_add_request_count",
		Help: "Общее количество объектов когда-либо помещенных в хранилище",
	},
)

var ObjectGetCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "object_get_request_count",
		Help: "Kоличество запрошенных объектов из хранилища",
	},
)

var ObjectDelCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "object_delete_request_count",
		Help: "Количество удаленных объектов из хранилища",
	},
)

var ObjectCountGauge = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "object_count_elements",
		Help: "Текущее количество объектов в хранилище",
	},
)
