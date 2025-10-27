package main

import (
	"net/http"
	"html/template"
)

type Page struct {
	Title string
	Content template.HTML
}

func InitAssets() {
	// Routes static assets like stylesheets and scripts

	styles := http.FileServer(http.Dir("web/styles/"))
	// scripts := http.FileServer(http.Dir("web/scripts/"))
	assets := http.FileServer(http.Dir("web/assets/"))

	http.Handle("/styles/", http.StripPrefix("/styles/", styles))
	// http.Handle("/scripts/", http.StripPrefix("/scripts/", scripts))
	http.Handle("/assets/", http.StripPrefix("/assets/", assets))
}

func InitIndex() {
	// Routes the main page and sets it up.

	page := Page{"Início", " "}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		template.Must(template.ParseFiles("web/index.html")).Execute(w, page)
	});
}

// func InitPosts()

// func InitArt()

func InitAbout() {
	page := Page{"Sobre", " "}

	http.HandleFunc("/sobre/", func(w http.ResponseWriter, r *http.Request) {
		template.Must(template.ParseFiles("web/index.html")).Execute(w, page)
	});

}

func main() {
	InitAssets()
	InitIndex()
	InitAbout()

	http.ListenAndServe(":8080", nil)
}
