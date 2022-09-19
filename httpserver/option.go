package httpserver

import "time"

func (s *Server) SetAddr(address string) {
	if address != "" {
		s.server.Addr = address
	}
}

func (s *Server) SetReadTimeout(timeout time.Duration) {
	s.server.ReadTimeout = timeout
}

func (s *Server) SetWriteTimeout(timeout time.Duration) {
	s.server.WriteTimeout = timeout
}

func (s *Server) SetIdleTimeout(timeout time.Duration) {
	s.server.IdleTimeout = timeout
}
