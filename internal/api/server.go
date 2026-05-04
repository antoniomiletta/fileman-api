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

func NewServer(authSvc *services.AuthService, fileSvc *services.FileService, folderSvc *services.FolderService, cfg config.ServerConfig) *Server {
	s := &Server{
		AuthSvc:   authSvc,
		FileSvc:   fileSvc,
		FolderSvc: folderSvc,
	}

	s.HttpServer = &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: s.Router(),
	}

	return s
}

func (s *Server) Start() {
	s.HttpServer.ListenAndServe()
}
