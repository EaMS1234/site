package handlers

import "net/http"


func Assets(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/assets/" + r.PathValue("file"))
}

func Styles(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/styles/" + r.PathValue("file"))
}

func Static(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "content/static/" + r.PathValue("file"))
}

