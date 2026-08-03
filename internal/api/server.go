package api

import (
	"fmt"
	"net/http"

	"github.com/antoniomiletta/fileman/config"
	"github.com/antoniomiletta/fileman/internal/api/transport"
	"github.com/antoniomiletta/fileman/internal/pkg/authenticator"
	"github.com/antoniomiletta/fileman/internal/service"
)

type Server struct {
	HttpServer *http.Server
	AuthSvc    *service.AuthService
	FileSvc    *service.FileService
	FolderSvc  *service.FolderService
	Authn      *authenticator.Authenticator
	Resp       *transport.Responder
}

type ServerDeps struct {
	Cfg       config.ServerConfig
	AuthSvc   *service.AuthService
	FolderSvc *service.FolderService
	FileSvc   *service.FileService
	Authn     *authenticator.Authenticator
	Resp      *transport.Responder
}

func NewServer(deps ServerDeps) *Server {
	s := &Server{
		AuthSvc:   deps.AuthSvc,
		FolderSvc: deps.FolderSvc,
		FileSvc:   deps.FileSvc,
		Authn:     deps.Authn,
		Resp:      deps.Resp,
	}

	s.HttpServer = &http.Server{
		Addr:    ":" + deps.Cfg.Port,
		Handler: s.Router(),
	}

	return s
}

func (s *Server) Serve(cfg config.ServerConfig) error {
	fmt.Printf("listening on %s...\n", cfg.Port)
	return s.HttpServer.ListenAndServe()
}
