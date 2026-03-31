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


type Post struct {
	FileName string        // The full name of an image file for pictures
	Name string            // The name of the .md file (post) or the image file (pictures) without the extension
	Title string           // The title that is shown to the user
	Desc string            // The description/summary of the content
	Alt string             // The alternative name/URL of the content in another language
	TimeString string      // Actual timestamp that is shown to the user
	Time time.Time         // Timestamp of the content (used only for comparing and sorting)
	Content template.HTML  // Parsed contents of an .md file. Empty if no .md file (in the case of a picture)
}


type Index struct {
	Posts []Post
	Pictures []Post
	Year string
}


// Persistent list of posts. Is dynamically updated during runtime.
var posts = make(map[string][]Post)

// Persistent list of pictures. Also updated.
var pictures = make(map[string][]Post)


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
		goldmark.WithRendererOptions(html.WithUnsafe(), html.WithXHTML()),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	)

	markdown.Convert(file, &buf)

	// Parsed Markdown into HTML code
	html := template.HTML(buf.String())  
	
	// Parsed metadata
	md := markdown.Parser().Parse(text.NewReader([]byte(file))).OwnerDocument().Meta()


	// The metadata allow you to customize the description, url, timestamp or title of a content.
	// It also allows you to link to the same content in another language

	data["alt"] = fmt.Sprint(md["alt"])      // URL in the alternative language
	data["desc"] = fmt.Sprint(md["desc"])    // custom summary/description
	data["time"] = fmt.Sprint(md["time"])    // custom timestamp
	data["title"] = fmt.Sprint(md["title"])  // title of the post or picture
	data["url"] = fmt.Sprint(md["url"])      // custom URL

	// Empty string if no description
	if data["desc"] == "<nil>" {
		data["desc"] = ""
	}

	// Limits the description to 128 characters and add three dots
	if len(data["desc"]) >= 128 {
		data["desc"] = data["desc"][:128] + "..."
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
				file_name := file.Name()
				canonical_name := file_name[:len(file_name)-3]

				content, data := GetHtml("content/posts/" + lang + "/" + file_name)
				tm := GetTime(data["time"])

				if data["title"] == "<nil>" {
					data["title"] = canonical_name
				}

				if data["url"] == "<nil>" {
					data["url"] = canonical_name
				}

				posts[lang] = append(posts[lang], Post{"", data["url"], data["title"], data["desc"], data["alt"], tm.Format("02/01/2006 - 15:04"), tm, content})
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

				file_name := file.Name()
				canonical_name := file_name[:len(file_name)-4]  // get rid of the .jpg or .png

				// Gets the data from the file
				content, data := GetHtml("content/pictures/" + lang + "/" + file_name + ".md")

				// Fixes the empty strings in case there's no markdown for the image
				if data == nil {
					data = make(map[string]string)

					file, err := os.Stat("content/pictures/" + file_name)
					if err != nil {panic(err)}
					data["time"] = file.ModTime().Format("02/01/2006 - 15:04")
					data["title"] = canonical_name
					data["url"] = canonical_name
				}

				if data["title"] == "<nil>" {
					data["title"] = canonical_name
				}

				if data["url"] == "<nil>" {
					data["url"] = canonical_name
				}

				tm := GetTime(data["time"])

				pictures[lang] = append(pictures[lang], Post{file_name, data["url"], data["title"], data["desc"], data["alt"], tm.Format("02/01/2006 - 15:04"), tm, content})
			}
		}

		// Sorts the list by time
		sort.Slice(pictures[lang], func(a, b int) bool {
			return pictures[lang][a].Time.After(pictures[lang][b].Time)
		})
	}
}


// Custom 404 page
func handle404(w http.ResponseWriter, r *http.Request) {
	lang := ""

	if strings.Contains(r.URL.Path, "/en/") || !strings.Contains(r.Header.Get("Accept-Language"), "pt") {
		lang = "en"
	}

	w.WriteHeader(404)
	template.Must(template.ParseFiles("web/" + lang + "/404.html")).Execute(w, r.URL.String())
}

