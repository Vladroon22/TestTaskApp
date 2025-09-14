package server

import (
	"context"
	"net/http"
	"time"
)

type HttpServer struct {
	srv *http.Server
}

func NewServer(addr string, handler http.Handler) *HttpServer {
	s := &http.Server{
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		Handler:      handler,
	}

	return &HttpServer{srv: s}
}
func (h *HttpServer) Start() error {
	return h.srv.ListenAndServe()
}

func (h *HttpServer) Shutdown(ctx context.Context) error {
	return h.srv.Shutdown(ctx)
}
