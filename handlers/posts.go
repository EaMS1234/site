package handlers

import (
	"html/template"
	"net/http"
	"regexp"
	"strings"
)


func Posts(w http.ResponseWriter, r *http.Request) {
	target := r.PathValue("post")

	lang := ""

	if r.URL.Path == "/en/blog/" + target {
		lang = "en"
	}

	// Fetches the correct post from the list
	for _, post := range posts[lang] {
		if post.Name == target {
	 		template.Must(template.ParseFiles("web/" + lang + "/post.html")).Execute(w, post)
			return
		} 
	}
	
	handle404(w, r)
}


// Handles posts
func Blog(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/artigos/" && r.URL.Path != "/en/blog/" {
		reg := regexp.MustCompile("^/artigos/([^/]+)|^/en/blog/([^/]+)")
		http.Redirect(w, r, reg.FindString(r.URL.Path), http.StatusMovedPermanently)
		return
	}
	
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

	var list []Post;

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
				blog_page.Years[year] = Index{append(blog_page.Years[year].Posts, post), []Post{}, year}
			}
		} else {
			// If search is empty, append everything
			blog_page.Years[year] = Index{append(blog_page.Years[year].Posts, post), []Post{}, year}
		}
	}

	if en {
		template.Must(template.ParseFiles("web/en/blog.html")).Execute(w, blog_page)
	} else {
		template.Must(template.ParseFiles("web/blog.html")).Execute(w, blog_page)
	}
}

