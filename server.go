package main

import (
	"log"
	"net/http"
	"os"
)

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) WriteHeader(statusCode int) {
	// Remember the status code for later logging.
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

type notesHandler struct {
	http.Handler
	contentType string
}

func (handler *notesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sw := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
	if handler.contentType != "" {
		sw.Header().Set("Content-Type", handler.contentType)
	}
	handler.Handler.ServeHTTP(sw, r)
	log.Printf("%s \"%s %s %s\" %d \"%s\"",
		r.RemoteAddr,
		r.Method,
		r.URL.String(),
		r.Proto,
		sw.statusCode,
		r.UserAgent())
}

func NewServer(listenAddress string, root *os.Root, contentType string) *http.Server {
	handler := &notesHandler{
		Handler: http.FileServerFS(
			dotFileHidingFs{root.FS()},
		),
		contentType: contentType,
	}
	mux := http.NewServeMux()
	mux.Handle("/", handler)
	return &http.Server{Addr: listenAddress, Handler: mux}
}
