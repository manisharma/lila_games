package server

import (
	"context"
	"lila_games/pkg/service"
	"log"
	"net/http"
	"time"
)

// Server represents game serevr
type Server struct {
	port    string
	svc     *service.Service
	httpSrv *http.Server
}

// New creates new instance of game server
func New(port string, svc *service.Service) *Server {
	if svc == nil {
		log.Fatal("nil game service")
		return nil
	}

	http.HandleFunc("/", svc.RootHandler)
	http.HandleFunc("/top-modes", svc.TopModesHandler)
	http.HandleFunc("/player", svc.PlayerHandler)

	return &Server{
		port: port,
		svc:  svc,
		httpSrv: &http.Server{
			Addr:         ":1234",
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			Handler:      http.DefaultServeMux,
		},
	}
}

// Start starts the game server
func (s *Server) Start() {
	log.Print("game service listening on port" + s.port)
	go func() {
		if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to listen, err: " + err.Error())
		}
	}()
}

// Stop stops the game server
func (s *Server) Stop(ctx context.Context) error {
	if err := s.svc.Close(); err != nil {
		return err
	}
	return s.httpSrv.Shutdown(ctx)
}
