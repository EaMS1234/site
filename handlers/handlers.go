package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)


type Page struct {
	Name string
	Title string
	Desc string
	TimeString string
	Time time.Time
	Content template.HTML
	Alt string
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
func GetTime(timestamp string) time.Time {
	t, err := time.Parse("02/01/2006 - 15:04", timestamp)
	if err != nil {
		panic(err)
	}

	return t
}


// Converts a markdown located at "Path" to valid html
func GetHtml(Path string) (template.HTML, map[string]string) {
	file, err := os.ReadFile(Path)
	if err != nil {return template.HTML(""), nil}

	var buf bytes.Buffer
	var data = make(map[string]string)

	markdown := goldmark.New(
		goldmark.WithExtensions(meta.New(meta.WithStoresInDocument(),),),
		goldmark.WithRendererOptions(html.WithUnsafe()),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	)

	markdown.Convert(file, &buf)

	// Parsed Markdown into HTML code
	html := template.HTML(buf.String())  
	
	// Parsed metadata
	md := markdown.Parser().Parse(text.NewReader([]byte(file))).OwnerDocument().Meta()

	data["alt"] = fmt.Sprint(md["alt"])
	data["desc"] = fmt.Sprint(md["desc"])
	data["time"] = fmt.Sprint(md["time"])
	data["title"] = fmt.Sprint(md["title"])

	// Empty string if no description
	if data["desc"] == "<nil>" {
		data["desc"] = ""
	}

	if data["title"] == "<nil>" {
		data["title"] = Path
	}

	// Limits the description to 128 characters
	if len(data["desc"]) >= 128 {
		data["desc"] = data["desc"][:128]
	}

	// If no alt, then root
	if data["alt"] == "<nil>" {
		data["alt"] = "/"
	}

	// If no time, then get the file timestamp
	if data["time"] == "<nil>" {
		file, err := os.Stat(Path)
		if err != nil {panic(err)}

		data["time"] = file.ModTime().Format("02/01/2006 - 15:04")
	}

	return html, data
}


// Returns a list of the files on "Path"
func GetPosts() {
	for _, lang := range []string{"", "en"} {
		posts[lang] = nil

		data, err := os.ReadDir("content/posts/" + lang + "/")
		if err != nil {panic(err)}

		for _, file := range data {
			if !file.IsDir() {

				content, data := GetHtml("content/posts/" + lang + "/" + file.Name())
				tm := GetTime(data["time"])

				posts[lang] = append(posts[lang], Page{file.Name()[:len(file.Name())-3], data["title"], data["desc"], tm.Format("02/01/2006 - 15:04"), tm, content, data["alt"]})
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
				content, data := GetHtml("content/pictures/" + lang + "/" + file.Name() + ".md")
				
				// Fixes the empty string in case there's no markdown for the image
				if data == nil {
					data = make(map[string]string)

					file, err := os.Stat("content/pictures/" + file.Name())
					if err != nil {panic(err)}
					data["time"] = file.ModTime().Format("02/01/2006 - 15:04")
				}

				tm := GetTime(data["time"])

				pictures[lang] = append(pictures[lang], Picture{file.Name()[:len(file.Name())-4], data["desc"], tm.Format("02/01/2006 - 15:04"), file.Name(), tm, content})
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
		txt, data := GetHtml("content/about.en.md")

		template.Must(template.ParseFiles("web/en/content.html")).Execute(w, Page{"", "About", "", "", time.Now(), txt, data["alt"]})
	} else {
		txt, data := GetHtml("content/about.md")

		template.Must(template.ParseFiles("web/content.html")).Execute(w, Page{"", "Sobre", "", "", time.Now(), txt, data["alt"]})
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

