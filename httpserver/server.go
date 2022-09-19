package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	_defaultAddr            = "127.0.0.1:80"
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultIdleTimeout     = 8 * time.Second
	_defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	server          *http.Server
	ShutdownTimeout time.Duration
	PreStartFunc    func() error
	PostStopFunc    func() error
}

func New(handler http.Handler) Server {
	httpServer := &http.Server{
		Addr:         _defaultAddr,
		Handler:      handler,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultReadTimeout,
		IdleTimeout:  _defaultIdleTimeout,
	}

	server := Server{
		server:          httpServer,
		ShutdownTimeout: _defaultShutdownTimeout,
	}

	server.PreStartFunc = func() error {
		fmt.Println("HTTP Server started on", httpServer.Addr)
		return nil
	}
	server.PostStopFunc = func() error {
		fmt.Println("HTTP Server stopped")
		return nil
	}

	return server
}

func (s Server) ListenAndServe() error {
	if err := s.PreStartFunc(); err != nil {
		return err
	}
	return s.server.ListenAndServe()
}

func (s Server) Shutdown(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	return s.PostStopFunc()
}
