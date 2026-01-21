package handlers

import (
	"net/http"
	"os"
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

func Static(w http.ResponseWriter, r *http.Request) {
	_, err := os.Stat("web/static/" + r.PathValue("file"))
	if err != nil {
		handle404(w, r)
		return
	}

	http.ServeFile(w, r, "content/static/" + r.PathValue("file"))
}

