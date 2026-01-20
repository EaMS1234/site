package handlers

import (
	"net/http"
	"os"
	"html/template"
	"strings"
)


// Handles posts
func InitPosts(content_dir string) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		en := (r.URL.Path == "/en/blog/")

		target := r.URL.Query().Get("a")
		search := r.URL.Query().Get("q")

		if target != "" {
			var file string

			if en {
				file = content_dir + "/posts/en/" + target + ".md"
			} else {
				file = content_dir + "/posts/" + target + ".md"
			}

			// Checks if the file exists
			_, err := os.Stat(file)
			if err != nil {
				if en {
					handle404(w, r, "en")
				} else {
					handle404(w, r, "")
				}
				return
			}

			content := GetHtml(file)
			tm := GetTime(file)

			if en {
				template.Must(template.ParseFiles("web/en/content.html")).Execute(w, Page{target, "", tm.Format("02/01/2006 - 15:04"), tm, content})
			} else {
				template.Must(template.ParseFiles("web/content.html")).Execute(w, Page{target, "", tm.Format("02/01/2006 - 15:04"), tm, content})
			}
		} else {
			var posts = struct{Title string; Years map[string]Index}{"Todos os artigos", make(map[string]Index)}

			if en {
				posts.Title = "Blog"
			}

			if search != "" {
				if en {
					posts.Title = "Results for \"" + search + "\""
				} else {
					posts.Title = "Resultados para \"" + search + "\""
				}
			}

			var list []Page;

			if en {
				list = GetPosts(content_dir + "/posts/en")
			} else {
				list = GetPosts(content_dir + "/posts")
			}

			for _, post := range list {
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

			if en {
				template.Must(template.ParseFiles("web/en/posts.html")).Execute(w, posts)
			} else {
				template.Must(template.ParseFiles("web/posts.html")).Execute(w, posts)
			}
		}
	} 

	http.HandleFunc("/artigos/", handler);
	http.HandleFunc("/en/blog/", handler);
}

