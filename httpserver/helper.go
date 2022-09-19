package httpserver

import (
	"context"
	"net/http"
)

func (s Server) Listen(ctx context.Context, addr string) error {
	shutdown := make(chan error, 1)
	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), s.ShutdownTimeout)
		defer cancel()

		shutdown <- s.Shutdown(ctx)
	}()

	s.SetAddr(addr)

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return <-shutdown
}
