package handlers

import (
	"net/http"
	"html/template"
	// "strings"
)


// Routes the main page and sets it up.
func InitIndex(w http.ResponseWriter, r *http.Request) {
	var index Index	

	en := (r.URL.Path == "/en/")

	// Updates the list of posts asynchronously
	go GetPosts()

	if en {
		index.Posts = posts["en"]
		index.Pictures = GetPictures("content/pictures/", "/en")
	} else {
		index.Posts = posts[""]
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

