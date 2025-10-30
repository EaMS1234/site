package main

import (
	"bytes"
	"html/template"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/yuin/goldmark"
)


type Page struct {
	Title string
	Desc string
	Time string
	Content template.HTML
}


type Index struct {
	Posts []Page
}


// Simple function for getting the timestamp of a file
func GetTime(File string) time.Time {
	file, err := os.Stat(File)
	if err != nil {panic(err)}

	return file.ModTime()
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


// Returns a list of the files on "Path"
func GetPosts(Path string) []Page {
	data, err := os.ReadDir(Path)
	if err != nil {panic(err)}

	var list []Page

	for _, file := range data {
		if !file.IsDir() {

			// Gets the timestamp for the file
			tm := GetTime(Path + "/" + file.Name())

			// Gets the first 128 characters as a description
			f, err := os.Open((Path + "/" + file.Name()))
			if err != nil {panic(err)}

			buf := make([]byte, 128)

			head, err := f.Read(buf)
			if err != nil {panic(err)}


			list = append(list, Page{file.Name()[:len(file.Name())-3], string(buf[:head]), tm.Format("1/2/2006 - 15:04"), ""})
		}
	}

	// Sorts the list by time
	sort.Slice(list, func(a, b int) bool {
		return list[a].Time > list[b].Time
	})

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
	var index Index

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		index.Posts = GetPosts("content/posts/")

		template.Must(template.ParseFiles("web/index.html")).Execute(w, index)
	})
}


// Routes the "about" page and sets it up.
func InitAbout() {
	page := Page{"Sobre", "", "", ""}

	page.Content = GetHtml("content/about.md")

	http.HandleFunc("/sobre/", func(w http.ResponseWriter, r *http.Request) {
		template.Must(template.ParseFiles("web/content.html")).Execute(w, page)
	});

}


// Handles posts
func InitPosts() {
	http.HandleFunc("/artigo", func(w http.ResponseWriter, r *http.Request) {
		target := r.URL.Query().Get("a")

		file := "content/posts/" + target + ".md"

		content := GetHtml(file)
		tm := GetTime(file).Format("02-01-2006 - 15:04")

		template.Must(template.ParseFiles("web/content.html")).Execute(w, Page{target, "", tm, content})
	});
}


func main() {
	InitAssets()
	InitIndex()
	InitAbout()
	InitPosts()

	http.ListenAndServe(":8080", nil)
}

