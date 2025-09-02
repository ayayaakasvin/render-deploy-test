package httpserver

import (
	"net/http"
	"sync"
	"time"

	"github.com/ayayaakasvin/lightmux"
	"github.com/render-test-server/internal/config"
	"github.com/sirupsen/logrus"
)

type ServerApp struct {
	server 	*http.Server

	lmux   	*lightmux.LightMux
	logger 	*logrus.Logger
	cfg 	*config.Config

	wg 		*sync.WaitGroup
}

func NewServerApp(cfg *config.Config, logger *logrus.Logger, wg *sync.WaitGroup) *ServerApp {
	serverapp := &ServerApp{}

	serverapp.server = &http.Server{}
	serverapp.lmux = lightmux.NewLightMux(serverapp.server)
	serverapp.logger = logger
	serverapp.cfg = cfg
	serverapp.wg = wg

	return serverapp
}

func (s *ServerApp) Run() {
	defer s.wg.Done()

	s.setupServer()

	s.setupLightMux()

	s.startServer()
}

func (s *ServerApp) startServer() {
	s.logger.Infof("Server has been started on port: %s", s.cfg.Address)
	s.logger.Infof("Available handlers:\n")

	s.lmux.PrintMiddlewareInfo()
	s.lmux.PrintRoutes()

	go func() {
		ticker := time.NewTicker(time.Minute * 5)
		for range ticker.C {
			s.logger.Info("Server is running...")
		}
	}()

	// RunTLS can be run when server is hosted on domain, acts as seperate service of file storing, for my project, id chose to encapsulate servers under one docker-compose and make nginx-gateaway for my api like auth, file, user service
	// if err := s.lmux.RunTLS(s.cfg.TLS.CertFile, s.cfg.TLS.KeyFile); err != nil {
	if err := s.lmux.Run(); err != nil {
		s.logger.Fatalf("Server exited with error: %v", err)
	}
}

// setuping server by pointer, so we dont have to return any value
func (s *ServerApp) setupServer() {
	if s.server == nil {
		// s.logger.Warn("Server is nil, creating a new server pointer")
		s.server = &http.Server{}
	}

	s.server.Addr = s.cfg.Address
	s.server.IdleTimeout = s.cfg.IdleTimeout
	s.server.ReadTimeout = s.cfg.Timeout
	s.server.WriteTimeout = s.cfg.Timeout

	s.logger.Info("Server has been set up")
}

func (s *ServerApp) setupLightMux() {
	s.lmux = lightmux.NewLightMux(s.server)

	s.lmux.NewRoute("/ping").Handle(http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
		w.WriteHeader(http.StatusOK)
	})

	s.lmux.NewRoute("/hello", func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			agent := r.UserAgent()
			method := r.Method
			url := r.URL.String()
			s.logger.Infof("Incoming request: IP=%s, Agent=%s, Method=%s, URL=%s", ip, agent, method, url)
			next(w, r)
		}
	}).Handle(http.MethodGet, HelloHandler)

	s.logger.Info("LightMux has been set up")
}
