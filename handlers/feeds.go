package handlers

import (
	"net/http"
	"text/template"
)

func Feeds(w http.ResponseWriter, r *http.Request) {
	GetPosts()
	GetPictures()

	if r.URL.Path == "/artigos/feed" {
		list := posts[""]

		for i := range list {
			list[i].TimeString = list[i].Time.Format("2006-01-02T15:04:05Z")
		}

		template.Must(template.ParseFiles("web/feeds/blog.xml")).Execute(w, struct{Update string; List []Post}{list[0].TimeString, list})
		return
	}

	if r.URL.Path == "/en/blog/feed" {
		list := posts["en"]

		for i := range list {
			list[i].TimeString = list[i].Time.Format("2006-01-02T15:04:05Z")
		}

		template.Must(template.ParseFiles("web/en/feeds/blog.xml")).Execute(w, struct{Update string; List []Post}{list[0].TimeString, list})
		return
	}
}

