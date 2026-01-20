package handlers

import (
	"html/template"
	"net/http"
	"os"
	"strings"
)


func Posts(w http.ResponseWriter, r *http.Request) {
	target := r.PathValue("post")
	file := "content/posts/" + target + ".md"

	en := (r.URL.Path == "/en/blog/" + target)

	if en {
		file = "content/posts/en/" + target + ".md"
	}

	// Checks if the file exists
	_, err := os.Stat(file)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	content := GetHtml(file)
	tm := GetTime(file)

	if en {
	 	template.Must(template.ParseFiles("web/en/content.html")).Execute(w, Page{target, "", tm.Format("02/01/2006 - 15:04"), tm, content})
	} else {
	 	template.Must(template.ParseFiles("web/content.html")).Execute(w, Page{target, "", tm.Format("02/01/2006 - 15:04"), tm, content})
	}
}


// Handles posts
func Blog(w http.ResponseWriter, r *http.Request) {
	en := (r.URL.Path == "/en/blog/")

	search := r.URL.Query().Get("q")

	// Struct that defines the page
	var blog_page = struct{Title string; Years map[string]Index}{"Todos os artigos", make(map[string]Index)}

	if en {
		blog_page.Title = "Blog"
	}

	if search != "" {
		if en {
			blog_page.Title = "Results for \"" + search + "\""
		} else {
			blog_page.Title = "Resultados para \"" + search + "\""
		}
	}

	var list []Page;

	// Updates the list of posts asynchronously
	go GetPosts()

	if en {
		list = posts["en"]
	} else {
		list = posts[""]
	}

	for _, post := range list {
		year := post.Time.Format("2006")
				
		if search != "" {
			// lowers everything to make it NOT case sensitive
			search = strings.ToLower(search)

			// If search is not empty, append only the matching content
			if strings.Contains(strings.ToLower(post.Title), search) || strings.Contains(strings.ToLower(post.Desc), search) || strings.Contains(post.Time.Format("2006"), search) {
				blog_page.Years[year] = Index{append(blog_page.Years[year].Posts, post), []Picture{}, year}
			}
		} else {
			// If search is empty, append everything
			blog_page.Years[year] = Index{append(blog_page.Years[year].Posts, post), []Picture{}, year}
		}
	}

	if en {
		template.Must(template.ParseFiles("web/en/posts.html")).Execute(w, blog_page)
	} else {
		template.Must(template.ParseFiles("web/posts.html")).Execute(w, blog_page)
	}
}

