package reconfigurableserver

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Server struct {
	server                   *http.Server
	mux                      *http.ServeMux
	mutex                    sync.Mutex
	busy                     bool
	ShutdownSignal           chan os.Signal
	listenaddr               string
	CertificateFile, KeyFile string
}

func (s *Server) limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.Busy() {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) Mux() *http.ServeMux {
	return s.mux
}

func (s *Server) GetNewServer() *http.Server {
	return &http.Server{
		Addr:    s.listenaddr,
		Handler: s.limit(s.mux),
	}
}

func (s *Server) Busy() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.busy
}

func NewServer(listenaddr, CertificateFile, KeyFile string) *Server {
	s := &Server{
		CertificateFile: CertificateFile,
		KeyFile:         KeyFile,
		listenaddr:      listenaddr,
		mux:             http.NewServeMux(),
		ShutdownSignal:  make(chan os.Signal, 1),
	}
	s.server = s.GetNewServer()
	return s
}

func (s *Server) Start() {
	if err := s.server.ListenAndServeTLS(s.CertificateFile, s.KeyFile); err != http.ErrServerClosed {
		log.Fatalln("Fatal Error", err)
	} else {
		log.Println("Warning: Server::Start() shutting down")
	}
}

func (s *Server) setBusy(busy bool) {
	s.mutex.Lock()
	s.mutex.Unlock()
	s.busy = busy
}

func (s *Server) RestartServer() {
	s.setBusy(true)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	s.server.Shutdown(ctx)
	s.server = s.GetNewServer()
	go s.Start()
	s.setBusy(false)
}
