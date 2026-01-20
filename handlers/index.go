package handlers

import (
	"net/http"
	"html/template"
	// "strings"
)


// Routes the main page and sets it up.
func InitIndex(w http.ResponseWriter, r *http.Request) {
	var index Index	

	// if r.URL.Path != "/en/" && r.URL.Path != "/" {
	// 	languages := r.Header.Get("Accept-Languages")
	// 	if strings.Contains(languages, "pt") || strings.Contains(languages, "pt-BR") {
	// 		handle404(w, r, "")
	// 	} else {
	// 		handle404(w, r, "en")
	// 	}
	// 	return
	// }

	en := (r.URL.Path == "/en/")

	if en {
		index.Posts = GetPosts("content/posts/en/")
		index.Pictures = GetPictures("content/pictures/", "/en")
	} else {
		index.Posts = GetPosts("content/posts/")
		index.Pictures = GetPictures("content/pictures/", "")
	}

	// Show at most 3 posts
	if len(index.Posts) >= 3 {
		index.Posts = index.Posts[:3]
	}

	// Show at most only one picture
	if len(index.Pictures) > 1 {
		index.Pictures = index.Pictures[0:1]
	}

	if en {
		template.Must(template.ParseFiles("web/en/index.html")).Execute(w, index)
	} else {
		template.Must(template.ParseFiles("web/index.html")).Execute(w, index)
	}
}

