package service

import "net/http"

type FileHandler struct{}

// ServeHTTP calls f(w, r).
func (f FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("file not found"))
	return
}
