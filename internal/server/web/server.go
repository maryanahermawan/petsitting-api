package web

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	Address  string
	basePath string
	router   *chi.Mux
}

func NewServer(address string) *Server {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	return &Server{
		Address: address,
		router:  r,
	}
}

func (server *Server) AddRoute(path string, fn http.Handler) {
	(*server).router.Mount(path, fn)
}

func (server *Server) Start(ctx context.Context) {
	s := &http.Server{Addr: server.Address, Handler: server.router}
	s.RegisterOnShutdown(func() {
		log.Infof("Shutting down server")
		ctx.Done()
	})

	err := s.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
		ctx.Done()
	}
}
