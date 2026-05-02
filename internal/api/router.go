package api

import (
	"net/http"

	"github.com/antoniomiletta/fileman/internal/api/handlers"
	"github.com/antoniomiletta/fileman/internal/api/middleware"
)

func (s *Server) Router() http.Handler {
	mux := http.DefaultServeMux
	authHandler := handlers.NewAuthHandler(s.AuthSvc)
	fileHandler := handlers.NewFileHandler(s.FileSvc)
	folderHandler := handlers.NewFolderHandler(s.FolderSvc)

	mux.HandleFunc("POST /auth/login", authHandler.Register)

	protected := middleware.Chain(middleware.Auth)

	mux.Handle("POST /files", protected(http.HandlerFunc(fileHandler.CreateFile)))
	mux.Handle("GET /files/{folderID}", protected(http.HandlerFunc(fileHandler.ListFromFolder)))

	mux.Handle("GET /folders", protected(http.HandlerFunc(folderHandler.CreateFolder)))
	mux.Handle("POST /folders/{folderID}", protected(http.HandlerFunc(folderHandler.ListChildren)))

	return mux
}
