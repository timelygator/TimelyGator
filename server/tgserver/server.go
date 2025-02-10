package tgserver

import (
    "fmt"
    "log"
    "net/http"
    "time"

    db "timelygator/server/database"
)

// Server holds fields for running your app's HTTP server.
type Server struct {
    Host    string
    Port    int
    Testing bool
    Router  http.Handler
    API     *ServerAPI
    Logger  *log.Logger
}

func NewServer(host string, port int, testing bool) (*Server, error) {
    logger := log.Default()

    // Initialize datastore
    dstore, err := db.NewDatastore(testing, nil, logger)
    if err != nil {
        return nil, fmt.Errorf("failed to init datastore: %w", err)
    }

    // Create the high-level ServerAPI
    api, err := NewServerAPI(dstore, testing, "1.0.0")
    if err != nil {
        return nil, fmt.Errorf("failed to init ServerAPI: %w", err)
    }

    return &Server{
        Host:    host,
        Port:    port,
        Testing: testing,
        Router:  nil,   // to be set externally
        API:     api,
        Logger:  logger,
    }, nil
}

// Start runs the server.
func (s *Server) Start() error {
    addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
    s.Logger.Printf("Starting server on %s (testing=%v)\n", addr, s.Testing)

    srv := &http.Server{
        Addr:         addr,
        Handler:      s.Router,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
    }

    return srv.ListenAndServe()
}
