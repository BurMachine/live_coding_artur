package transport

import (
	"context"
	"log"
	"net/http"
	"time"

	"awesomeProject/internal/config"
	"awesomeProject/internal/transport/handlers"
	"github.com/gorilla/mux"
)

type Server struct {
	cfg        *config.Config
	httpServer *http.Server
}

func New(cfg *config.Config) *Server {
	server := http.Server{
		Addr:              ":8888",
		ReadHeaderTimeout: 30 * time.Second,
	}

	return &Server{
		cfg:        cfg,
		httpServer: &server,
	}
}

func (s *Server) Start(ctx context.Context, handlers *handlers.Handlers) {
	router := setupRoutes(handlers)
	s.httpServer.Handler = router

	go func(ctx context.Context) {
		log.Fatal(s.httpServer.ListenAndServe())
	}(ctx)
	select {
	case <-ctx.Done():
		log.Println("shutting down http server")
		s.httpServer.Shutdown(context.Background())
	}
}

func setupRoutes(handlers *handlers.Handlers) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/get_count/{pageID}", handlers.GetCountReq).Methods(http.MethodGet)
	r.HandleFunc("/increment/{pageID}", handlers.IncrementReq).Methods(http.MethodPost)

	return r
}
