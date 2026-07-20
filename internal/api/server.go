package api

import (
	"net/http"

	"github.com/antoniomiletta/fileman/config"
	"github.com/antoniomiletta/fileman/internal/pkg/authenticator"
	"github.com/antoniomiletta/fileman/internal/service"
)

type Server struct {
	HttpServer    *http.Server
	AuthSvc       *service.AuthService
	FileSvc       *service.FileService
	FolderSvc     *service.FolderService
	Authenticator *authenticator.Authenticator
}

type ServerDeps struct {
	Cfg           config.ServerConfig
	AuthSvc       *service.AuthService
	FolderSvc     *service.FolderService
	FileSvc       *service.FileService
	Authenticator *authenticator.Authenticator
}

func NewServer(deps ServerDeps) *Server {
	s := &Server{
		AuthSvc:       deps.AuthSvc,
		FolderSvc:     deps.FolderSvc,
		FileSvc:       deps.FileSvc,
		Authenticator: deps.Authenticator,
	}

	s.HttpServer = &http.Server{
		Addr:    ":" + deps.Cfg.Port,
		Handler: s.Router(),
	}

	return s
}

func (s *Server) Start() {
	s.HttpServer.ListenAndServe()
}
