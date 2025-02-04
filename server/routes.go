package server

import (
	"github.com/alexchebotarsky/heatpump-api/server/handler"
	"github.com/alexchebotarsky/heatpump-api/server/middleware"
	chi "github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *Server) setupRoutes() {
	s.Router.Get("/_healthz", handler.Health)
	s.Router.Handle("/metrics", promhttp.Handler())

	s.Router.Route(v1API, func(r chi.Router) {
		r.Use(middleware.Metrics)

		r.Get("/heatpump/state", handler.GetHeatpumpState(s.Clients.Database))
		r.Post("/heatpump/state", handler.UpdateHeatpumpState(s.Clients.Database, s.Clients.PubSub))
	})
}

const v1API = "/api/v1"
