package main

import (
	"bytes"
	"html/template"
	"net/http"
	"os"

	"github.com/yuin/goldmark"
)


type Page struct {
	Title string
	Content template.HTML
}


// Converts a markdown located at "Path" to valid html
func GetHtml(Path string) template.HTML {
	file, err := os.ReadFile(Path)
	if err != nil {panic(err)}

	var buf bytes.Buffer
	if err := goldmark.Convert(file, &buf); err != nil {panic(err)}

	html := template.HTML(buf.String())

	return html
}


// Routes static assets like stylesheets and scripts
func InitAssets() {
	styles := http.FileServer(http.Dir("web/styles/"))
	// scripts := http.FileServer(http.Dir("web/scripts/"))
	assets := http.FileServer(http.Dir("web/assets/"))

	http.Handle("/styles/", http.StripPrefix("/styles/", styles))
	// http.Handle("/scripts/", http.StripPrefix("/scripts/", scripts))
	http.Handle("/assets/", http.StripPrefix("/assets/", assets))
}


// Routes the main page and sets it up.
func InitIndex() {
	page := Page{"Início", " "}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		template.Must(template.ParseFiles("web/index.html")).Execute(w, page)
	});
}


// Routes the "about" page and sets it up.
func InitAbout() {
	page := Page{"Sobre", " "}

	page.Content = GetHtml("content/about.md")

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

