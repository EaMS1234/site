package handlers

import (
	"net/http"
	"os"
	"strings"
)


func Assets(w http.ResponseWriter, r *http.Request) {
	_, err := os.Stat("web/assets/" + r.PathValue("file"))
	if err != nil {
		handle404(w, r)
		return
	}

	http.ServeFile(w, r, "web/assets/" + r.PathValue("file"))
}

func Styles(w http.ResponseWriter, r *http.Request) {
	_, err := os.Stat("web/styles/" + r.PathValue("file"))
	if err != nil {
		handle404(w, r)
		return
	}

	http.ServeFile(w, r, "web/styles/" + r.PathValue("file"))
}

func Scripts(w http.ResponseWriter, r *http.Request) {
	_, err := os.Stat("web/scripts/" + r.PathValue("file"))
	if err != nil {
		handle404(w, r)
		return
	}

	http.ServeFile(w, r, "web/scripts/" + r.PathValue("file"))
}

func Static(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/static/")

	_, err := os.Stat("content/static/" + path)
	if err != nil || path == "" {
		handle404(w, r)
		return
	}

	http.StripPrefix("/static/", http.FileServer(http.Dir("content/static"))).ServeHTTP(w, r)
}

