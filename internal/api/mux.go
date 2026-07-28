package api

import (
	"net/http"

	"github.com/antoniomiletta/fileman/internal/api/handler"
	"github.com/antoniomiletta/fileman/internal/api/middleware"
)

func (s *Server) Router() http.Handler {
	mux := http.DefaultServeMux

	authHandler := handler.NewAuthHandler(s.AuthSvc)
	fileHandler := handler.NewFileHandler(s.FileSvc)
	folderHandler := handler.NewFolderHandler(s.FolderSvc)

	mux.HandleFunc("POST /auth/signup", authHandler.SignUp)
	mux.HandleFunc("POST /auth/login", authHandler.Login)

	mw := middleware.New(s.Authenticator)
	protected := middleware.Chain(mw.Auth)

	mux.Handle("POST /folders", protected(http.HandlerFunc(folderHandler.Create)))
	mux.Handle("GET /folders/{id}", protected(http.HandlerFunc(folderHandler.ListChildren)))
	mux.Handle("PATCH /folders/{id}/move", protected(http.HandlerFunc(folderHandler.Move)))
	mux.Handle("DELETE /folders/{id}", protected(http.HandlerFunc(folderHandler.Delete)))

	mux.Handle("POST /files", protected(http.HandlerFunc(fileHandler.Create)))
	mux.Handle("PATCH /files/{id}/move", protected(http.HandlerFunc(fileHandler.Move)))
	mux.Handle("DELETE /files/{id}", protected(http.HandlerFunc(fileHandler.Delete)))

	return mux
}
