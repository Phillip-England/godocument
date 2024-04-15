package handler

import (
	"net/http"
	"path/filepath"
)

func ServeFavicon(w http.ResponseWriter, r *http.Request) {
	filePath := "favicon.ico"
	fullPath := filepath.Join(".", "static", filePath)
	http.ServeFile(w, r, fullPath)
}

func ServeStaticFiles(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[len("/static/"):]
	fullPath := filepath.Join(".", "static", filePath)
	http.ServeFile(w, r, fullPath)
}

func ServeOutFiles(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[len("/"):]
	fullPath := filepath.Join(".", "out", filePath)
	http.ServeFile(w, r, fullPath)
}
