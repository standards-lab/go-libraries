package web

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	http *http.Server
	errs chan error

	mu       sync.Mutex
	listener net.Listener
}

func NewServer(cfg Config, handler http.Handler) *Server {
	return &Server{
		http: &http.Server{
			Addr:              cfg.Addr(),
			Handler:           handler,
			ReadTimeout:       time.Duration(cfg.ReadTimeout),
			ReadHeaderTimeout: time.Duration(cfg.ReadHeaderTimeout),
			WriteTimeout:      time.Duration(cfg.WriteTimeout),
			IdleTimeout:       time.Duration(cfg.IdleTimeout),
		},
		errs: make(chan error, 1),
	}
}

func (s *Server) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.listener != nil {
		return errors.New("server already started")
	}

	listener, err := net.Listen("tcp", s.http.Addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", s.http.Addr, err)
	}
	s.listener = listener

	go func() {
		defer close(s.errs)
		if err := s.http.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.errs <- fmt.Errorf("serve %s: %w", listener.Addr(), err)
		}
	}()

	return nil
}

func (s *Server) Addr() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.listener != nil {
		return s.listener.Addr().String()
	}
	return s.http.Addr
}

func (s *Server) Err() <-chan error {
	return s.errs
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
