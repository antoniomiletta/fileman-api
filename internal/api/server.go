package api

import (
	"net/http"

	"github.com/antoniomiletta/fileman/config"
	"github.com/antoniomiletta/fileman/internal/services"
)

type Server struct {
	HttpServer *http.Server
	AuthSvc    *services.AuthService
	FileSvc    *services.FileService
	FolderSvc  *services.FolderService
}

type NewServerInput struct {
	AuthSvc   *services.AuthService
	FolderSvc *services.FolderService
	FileSvc   *services.FileService
	Cfg       config.ServerConfig
}

func NewServer(input NewServerInput) *Server {
	s := &Server{
		AuthSvc:   input.AuthSvc,
		FolderSvc: input.FolderSvc,
		FileSvc:   input.FileSvc,
	}

	s.HttpServer = &http.Server{
		Addr:    ":" + input.Cfg.Port,
		Handler: s.Router(),
	}

	return s
}

func (s *Server) Start() {
	s.HttpServer.ListenAndServe()
}
