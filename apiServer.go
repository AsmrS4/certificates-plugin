package certificatesplugin

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (server *Server) Run(port string) error {
	server.httpServer = &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return server.httpServer.ListenAndServe()
}

func (s *Server) ServerShutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
