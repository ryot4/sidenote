package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

type dotFileHidingFs struct {
	http.FileSystem
}

func (fs dotFileHidingFs) Open(name string) (http.File, error) {
	for _, part := range strings.Split(name, "/") {
		if strings.HasPrefix(part, ".") {
			return nil, os.ErrNotExist
		}
	}

	file, err := fs.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}

	return dotFileHidingFile{file}, nil
}

type dotFileHidingFile struct {
	http.File
}

func (dir dotFileHidingFile) Readdir(count int) (filtered []os.FileInfo, err error) {
	files, err := dir.File.Readdir(count)
	for _, f := range files {
		if !strings.HasPrefix(f.Name(), ".") {
			filtered = append(filtered, f)
		}
	}
	return
}

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) WriteHeader(statusCode int) {
	// Remember the status code for later logging.
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

type requestLoggingHandler struct {
	http.Handler
}

func (handler *requestLoggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sw := &statusResponseWriter{ResponseWriter: w}
	handler.Handler.ServeHTTP(sw, r)
	log.Printf("%s \"%s %s %s\" %d \"%s\"",
		r.RemoteAddr,
		r.Method,
		r.URL.String(),
		r.Proto,
		sw.statusCode,
		r.UserAgent())
}

func NewServer(listenAddress, documentRoot string) *http.Server {
	handler := &requestLoggingHandler{
		http.FileServer(
			dotFileHidingFs{http.Dir(documentRoot)},
		),
	}
	mux := http.NewServeMux()
	mux.Handle("/", handler)
	return &http.Server{Addr: listenAddress, Handler: mux}
}
