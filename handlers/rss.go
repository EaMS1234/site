package handlers

import (
	"net/http"
	"text/template"
)

func Feed(w http.ResponseWriter, r *http.Request) {
	template.Must(template.ParseFiles("web/feed.xml")).Execute(w, struct{Posts []Page}{posts[""]})
}

