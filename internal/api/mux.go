package api

import (
	"net/http"

	"github.com/antoniomiletta/fileman/internal/api/handler"
	"github.com/antoniomiletta/fileman/internal/api/middleware"
)

func (s *Server) Router() http.Handler {
	mux := http.DefaultServeMux

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	authHandler := handler.NewAuthHandler(s.AuthSvc, s.Resp)
	fileHandler := handler.NewFileHandler(s.FileSvc, s.Resp)
	folderHandler := handler.NewFolderHandler(s.FolderSvc, s.Resp)

	mux.HandleFunc("POST /auth/signup", authHandler.SignUp)
	mux.HandleFunc("POST /auth/login", authHandler.Login)

	mw := middleware.New(s.Authn, s.Resp)
	protected := middleware.Chain(mw.Auth)

	mux.Handle("POST /folders", protected(http.HandlerFunc(folderHandler.Create)))
	mux.Handle("GET /folders/{id}", protected(http.HandlerFunc(folderHandler.ListChildren)))
	mux.Handle("PATCH /folders/{id}/move", protected(http.HandlerFunc(folderHandler.Move)))
	mux.Handle("PATCH /folders/{id}/rename", protected(http.HandlerFunc(folderHandler.Rename)))
	mux.Handle("DELETE /folders/{id}", protected(http.HandlerFunc(folderHandler.Delete)))

	mux.Handle("POST /files", protected(http.HandlerFunc(fileHandler.Create)))
	mux.Handle("PATCH /files/{id}/move", protected(http.HandlerFunc(fileHandler.Move)))
	mux.Handle("PATCH /files/{id}/rename", protected(http.HandlerFunc(fileHandler.Rename)))
	mux.Handle("DELETE /files/{id}", protected(http.HandlerFunc(fileHandler.Delete)))
	mux.Handle("GET /files/{id}/download", protected(http.HandlerFunc(fileHandler.Download)))

	return mux
}
