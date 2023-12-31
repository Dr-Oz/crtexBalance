package main

import (
	"context"
	"crtexBalance/internal/config"
	"log"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	conf       *config.Config
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:         port,
		Handler:      handler,
		ReadTimeout:  time.Duration(s.conf.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.conf.WriteTimeout) * time.Second,
	}

	log.Println("сервер запущен")

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
