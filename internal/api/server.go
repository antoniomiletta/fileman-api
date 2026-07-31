package api

import (
	"net/http"

	"github.com/antoniomiletta/fileman/config"
	"github.com/antoniomiletta/fileman/internal/pkg/authenticator"
	"github.com/antoniomiletta/fileman/internal/service"
)

type Server struct {
	HttpServer *http.Server
	AuthSvc    *service.AuthService
	FileSvc    *service.FileService
	FolderSvc  *service.FolderService
	Authn      *authenticator.Authenticator
}

type ServerDeps struct {
	Cfg       config.ServerConfig
	AuthSvc   *service.AuthService
	FolderSvc *service.FolderService
	FileSvc   *service.FileService
	Authn     *authenticator.Authenticator
}

func NewServer(deps ServerDeps) *Server {
	s := &Server{
		AuthSvc:   deps.AuthSvc,
		FolderSvc: deps.FolderSvc,
		FileSvc:   deps.FileSvc,
		Authn:     deps.Authn,
	}

	s.HttpServer = &http.Server{
		Addr:    ":" + deps.Cfg.Port,
		Handler: s.Router(),
	}

	return s
}

func (s *Server) Serve() error {
	return s.HttpServer.ListenAndServe()
}
