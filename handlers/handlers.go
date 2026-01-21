package handlers

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
	"github.com/yuin/goldmark/renderer/html"
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


// Persistent list of posts. Is dynamically updated during runtime.
var posts = make(map[string][]Page)

// Persistent list of pictures. Also updated.
var pictures = make(map[string][]Picture)


// Simple function for getting the timestamp of a file
func GetTime(File string) time.Time {
	file, err := os.Stat(File)
	if err != nil {panic(err)}

	return file.ModTime()
}


// Converts a markdown located at "Path" to valid html
func GetHtml(Path string) template.HTML {
	file, err := os.ReadFile(Path)
	if err != nil {return template.HTML("")}

	var buf bytes.Buffer
	if err := goldmark.New(goldmark.WithRendererOptions(html.WithUnsafe())).Convert(file, &buf); err != nil {panic(err)}

	html := template.HTML(buf.String())

	return html
}


// Returns a list of the files on "Path"
func GetPosts() {
	for _, lang := range []string{"", "en"} {
		posts[lang] = nil

		data, err := os.ReadDir("content/posts/" + lang + "/")
		if err != nil {panic(err)}

		for _, file := range data {
			if !file.IsDir() {

				// Gets the timestamp for the file
				tm := GetTime("content/posts/" + lang + "/" + file.Name())

				// Gets the first 128 characters as a description
				f, err := os.Open(("content/posts/" + lang + "/" + file.Name()))
				if err != nil {panic(err)}

				buf := make([]byte, 128)

				head, err := f.Read(buf)
				if err != nil {panic(err)}


				posts[lang] = append(posts[lang], Page{file.Name()[:len(file.Name())-3], string(buf[:head]), tm.Format("02/01/2006 - 15:04"), tm, ""})
			}
		}

		// Sorts the list by time
		sort.Slice(posts[lang], func(a, b int) bool {
			return posts[lang][a].Time.After(posts[lang][b].Time)
		})
	}
}


func GetPictures() {
	for _, lang := range []string{"", "en"} {
		pictures[lang] = nil

		data, err := os.ReadDir("content/pictures/")
		if err != nil {panic(err)}

		// Regex for files ending in .png or .jpg
		re := regexp.MustCompile(`.*\.(jpg|png)$`)

		for _, file := range data {
			if !file.IsDir() && re.MatchString(file.Name()) {

				// Gets the timestamp for the file
				tm := GetTime("content/pictures/" + file.Name())

				// Gets the first 128 characters as a description
				f, err := os.Open(("content/pictures/" + lang + "/" + file.Name() + ".md"))
				if err != nil {
					pictures[lang] = append(pictures[lang], Picture{file.Name()[:len(file.Name())-4], "", tm.Format("02/01/2006 - 15:04"), file.Name(), tm, ""})
				} else {
					buf := make([]byte, 128)

					head, err := f.Read(buf)
					if err != nil {panic(err)}

					pictures[lang] = append(pictures[lang], Picture{file.Name()[:len(file.Name())-4], string(buf[:head]), tm.Format("02/01/2006 - 15:04"), file.Name(), tm, ""})
				}
			}
		}

		// Sorts the list by time
		sort.Slice(pictures[lang], func(a, b int) bool {
			return pictures[lang][a].Time.After(pictures[lang][b].Time)
		})
	}
}


// Routes the "about" page and sets it up.
func About(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/en/about/" {
		template.Must(template.ParseFiles("web/en/content.html")).Execute(w, Page{"About", "", "", time.Now(), GetHtml("content/about.en.md")})
	} else {
		template.Must(template.ParseFiles("web/content.html")).Execute(w, Page{"Sobre", "", "", time.Now(), GetHtml("content/about.md")})
	}
}


func handle404(w http.ResponseWriter, r *http.Request) {
	lang := ""

	if strings.Contains(r.URL.Path, "/en/") || !strings.Contains(r.Header.Get("Accept-Language"), "pt") {
		lang = "en"
	}

	w.WriteHeader(404)
	template.Must(template.ParseFiles("web/" + lang + "/404.html")).Execute(w, r.URL.String())
}

