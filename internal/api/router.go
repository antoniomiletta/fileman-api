package api

import (
	"net/http"

	"github.com/antoniomiletta/fileman/internal/api/handlers"
	"github.com/antoniomiletta/fileman/internal/api/middlewares"
)

func (s *Server) Router() http.Handler {
	mux := http.DefaultServeMux

	authHandler := handlers.NewAuthHandler(s.AuthSvc)
	fileHandler := handlers.NewFileHandler(s.FileSvc)
	folderHandler := handlers.NewFolderHandler(s.FolderSvc)

	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)

	protected := middlewares.Chain(middlewares.Auth)

	mux.Handle("POST /folders", protected(http.HandlerFunc(folderHandler.Create)))
	mux.Handle("GET /folders/{id}", protected(http.HandlerFunc(folderHandler.ListChildren)))
	mux.Handle("PATCH /folders/{id}/move", protected(http.HandlerFunc(folderHandler.Move)))
	mux.Handle("DELETE /folders/{id}", protected(http.HandlerFunc(folderHandler.Delete)))

	mux.Handle("POST /files", protected(http.HandlerFunc(fileHandler.Create)))
	mux.Handle("PATCH /files/{id}/move", protected(http.HandlerFunc(fileHandler.Move)))
	mux.Handle("DELETE /files/{id}", protected(http.HandlerFunc(fileHandler.Delete)))

	return mux
}
