package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/yuin/goldmark"
)


type Page struct {
	Title string
	Content template.HTML
	Time time.Time
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


func GetPosts(Path string) []Page {
	data, err := os.ReadDir(Path)
	if err != nil {panic(err)}

	var list []Page

	for _, file := range data {
		if !file.IsDir() {

			// Gets the timestamp for the file
			in, err := file.Info()
			if err != nil {panic(err)}
			tm := time.Unix(in.Sys().(*syscall.Stat_t).Ctim.Sec, 0)

			list = append(list, Page{file.Name()[:len(file.Name())-3], "", tm})
		}
	}

	return list
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
	page := Page{"Início", " ", time.Now()}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		list := GetPosts("content/posts/")

		page.Content = template.HTML("")

		for _, post := range list {
			page.Content += template.HTML(fmt.Sprintf(
				"<a href=\"/artigo?a=%v\"><h1>%v</h1></a>%v",
				post.Title,
				post.Title,
				post.Time.Format("02/01/2006"),
			));
		}

		template.Must(template.ParseFiles("web/index.html")).Execute(w, page)
	});
}


// Routes the "about" page and sets it up.
func InitAbout() {
	page := Page{"Sobre", " ", time.Now()}

	page.Content = GetHtml("content/about.md")

	http.HandleFunc("/sobre/", func(w http.ResponseWriter, r *http.Request) {
		template.Must(template.ParseFiles("web/index.html")).Execute(w, page)
	});

}


func InitPosts() {
	http.HandleFunc("/artigo", func(w http.ResponseWriter, r *http.Request) {
		target := r.URL.Query().Get("a")

		content := GetHtml(fmt.Sprintf("content/posts/%v.md", target))

		template.Must(template.ParseFiles("web/index.html")).Execute(w, Page{target, content, time.Now()})
	});
}


func main() {
	InitAssets()
	InitIndex()
	InitAbout()
	InitPosts()

	http.ListenAndServe(":8080", nil)
}

