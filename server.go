package yagaw

import (
	"fmt"
	"net/http"

	"github.com/pho3b/tiny-logger/logs"
	"github.com/pho3b/tiny-logger/logs/log_level"
)

var Log *logs.Logger = InitLogger(log_level.ErrorLvlName)

func InitLogger(logLevel log_level.LogLvlName) *logs.Logger {
	return logs.NewLogger().
		SetLogLvl(logLevel).
		ShowLogLevel(false).
		EnableColors(true).
		AddTime(true).
		AddDate(true)
}

type Server struct {
	address string
	port    int
	server  *http.Server
	router  *Router
}

func (s *Server) Run() {
	s.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.address, s.port),
		Handler: s.router,
	}

	Log.Debug(fmt.Sprintf("Starting server on address `%s:%d`", s.address, s.port))
	err := s.server.ListenAndServe()
	if err != nil {
		Log.Error(err)
	}
}

func (s *Server) GetRouter() *Router {
	return s.router
}

func NewServer(addr string, port int) *Server {
	return &Server{
		address: addr,
		port:    port,
		router:  NewRouter(),
	}
}
