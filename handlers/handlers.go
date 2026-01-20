package handlers

import (
	"bytes"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"sort"
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


func GetPictures(Path string, lang string) []Picture {
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
			f, err := os.Open((Path + lang + "/" + file.Name() + ".md"))
			if err != nil {
				list = append(list, Picture{file.Name()[:len(file.Name())-4], "", tm.Format("02/01/2006 - 15:04"), file.Name(), tm, ""})
			} else {
				buf := make([]byte, 128)

				head, err := f.Read(buf)
				if err != nil {panic(err)}

				list = append(list, Picture{file.Name()[:len(file.Name())-4], string(buf[:head]), tm.Format("02/01/2006 - 15:04"), file.Name(), tm, ""})
			}
		}
	}

	// Sorts the list by time
	sort.Slice(list, func(a, b int) bool {
		return list[a].Time.After(list[b].Time)
	})

	return list
}


// Routes the "about" page and sets it up.
func About(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/en/about/" {
		template.Must(template.ParseFiles("web/en/content.html")).Execute(w, Page{"About", "", "", time.Now(), GetHtml("content/about.en.md")})
	} else {
		template.Must(template.ParseFiles("web/content.html")).Execute(w, Page{"Sobre", "", "", time.Now(), GetHtml("content/about.md")})
	}
}


func handle404(w http.ResponseWriter, r *http.Request, lang string) {
	w.WriteHeader(404)
	template.Must(template.ParseFiles("web/" + lang + "/404.html")).Execute(w, r.URL.String())
}


