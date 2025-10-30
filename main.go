package main

import (
	"bytes"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
)


type Page struct {
	Title string
	Desc string
	TimeString string
	Time time.Time
	Content template.HTML
}


type Picture struct {
	Title string
	Summary string
	TimeString string
	FileName string
	Time time.Time
	Description template.HTML
}


type Index struct {
	Posts []Page
	Pictures []Picture
	Year string
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


			list = append(list, Page{file.Name()[:len(file.Name())-3], string(buf[:head]), tm.Format("02/01/2006 - 15:04"), tm, ""})
		}
	}

	// Sorts the list by time
	sort.Slice(list, func(a, b int) bool {
		return list[a].Time.After(list[b].Time)
	})

	return list
}


func GetPictures(Path string) []Picture {
	data, err := os.ReadDir(Path)
	if err != nil {panic(err)}

	var list []Picture

	// Regex for files ending in .png or .jpg
	re := regexp.MustCompile(`.*\.(jpg|png)$`)

	for _, file := range data {
		if !file.IsDir() && re.MatchString(file.Name()) {

			// Gets the timestamp for the file
			tm := GetTime(Path + "/" + file.Name())

			// Gets the first 128 characters as a description
			f, err := os.Open((Path + "/" + file.Name() + ".md"))
			if err != nil {panic(err)}

			buf := make([]byte, 128)

			head, err := f.Read(buf)
			if err != nil {panic(err)}

			list = append(list, Picture{file.Name()[:len(file.Name())-4], string(buf[:head]), tm.Format("02/01/2006 - 15:04"), file.Name(), tm, ""})
		}
	}

	// Sorts the list by time
	sort.Slice(list, func(a, b int) bool {
		return list[a].Time.After(list[b].Time)
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

	// Content
	pictures := http.FileServer(http.Dir("content/pictures"))
	http.Handle("/pictures/", http.StripPrefix("/pictures/", pictures))
}


// Routes the main page and sets it up.
func InitIndex() {
	var index Index

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		index.Posts = GetPosts("content/posts/")
		index.Pictures = GetPictures("content/pictures/")

		template.Must(template.ParseFiles("web/index.html")).Execute(w, index)
	})
}


// Routes the "about" page and sets it up.
func InitAbout() {
	page := Page{"Sobre", "", "", time.Now(), ""}

	page.Content = GetHtml("content/about.md")

	http.HandleFunc("/sobre/", func(w http.ResponseWriter, r *http.Request) {
		template.Must(template.ParseFiles("web/content.html")).Execute(w, page)
	});
}


// Handles posts
func InitPosts() {
	http.HandleFunc("/artigos/", func(w http.ResponseWriter, r *http.Request) {
		target := r.URL.Query().Get("a")
		search := r.URL.Query().Get("q")

		if target != "" {
			file := "content/posts/" + target + ".md"

			content := GetHtml(file)
			tm := GetTime(file)

			template.Must(template.ParseFiles("web/content.html")).Execute(w, Page{target, "", tm.Format("02/01/2006 - 15:04"), tm, content})
		} else {
			var posts = struct{Title string; Years map[string]Index}{"Todos os artigos", make(map[string]Index)}

			if search != "" {
				posts.Title = "Resultados para \"" + search + "\""
			}

			for _, post := range GetPosts("content/posts") {
				year := post.Time.Format("2006")
				
				if search != "" {

					// If search is not empty, append only the matching content
					if strings.Contains(post.Title, search) || strings.Contains(post.Desc, search) || strings.Contains(post.Time.Format("2006"), search) {
						posts.Years[year] = Index{append(posts.Years[year].Posts, post), []Picture{}, year}
					}
				} else {
					// If search is empty, append everything
					posts.Years[year] = Index{append(posts.Years[year].Posts, post), []Picture{}, year}
				}
			}

			template.Must(template.ParseFiles("web/posts.html")).Execute(w, posts)
		}
	});
}


func InitGallery() {
	http.HandleFunc("/galeria", func(w http.ResponseWriter, r *http.Request) {
		
	});
}


func main() {
	InitAssets()
	InitIndex()
	InitAbout()
	InitPosts()
	InitGallery()

	http.ListenAndServe(":8080", nil)
}

