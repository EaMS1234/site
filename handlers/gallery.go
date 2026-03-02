package handlers

import (
	"net/http"
	"os"
	"html/template"
	"strings"
)


func Images(w http.ResponseWriter, r *http.Request) {
	_, err := os.Stat("content/pictures/" + r.PathValue("file"))
	if err != nil {
		handle404(w, r)
		return
	}

	http.ServeFile(w, r, "content/pictures/" + r.PathValue("file"))
}


func Pictures(w http.ResponseWriter, r *http.Request) {
	image := r.PathValue("pic")

	lang := ""

	if r.URL.Path == "/en/pictures/" + image {
		lang = "en"
	}
		
	for _, pic := range(pictures[lang]) {
		if pic.Name == image {
			template.Must(template.ParseFiles("web/" + lang + "/picture.html")).Execute(w, pic)
			return
		}
	}

	handle404(w, r)
}


func Gallery(w http.ResponseWriter, r *http.Request) {
	en := (r.URL.Path == "/en/pictures/")

	search := r.URL.Query().Get("q")


	var gallery = struct{Title string; Years map[string]Index}{"Galeria", make(map[string]Index)}

	if search != "" {
		if en {
			gallery.Title = "Results for \"" + search + "\""
		} else {
			gallery.Title = "Resultados para \"" + search + "\""
		}
	} else if en {
		gallery.Title = "Pictures"
	}

	go GetPictures()

	var list []Post

	if en {
		list = pictures["en"]
	} else {
		list = pictures[""]
	}

	for _, pic := range list {
		year := pic.Time.Format("2006")

		if search != "" {
			search := strings.ToLower(search)

			// If search is not empty, append only the matching content
			if strings.Contains(strings.ToLower(pic.Title), search) || strings.Contains(strings.ToLower(pic.Desc), search) || strings.Contains(pic.Time.Format("2006"), search) {
				gallery.Years[year] = Index{[]Post{}, append(gallery.Years[year].Pictures, pic), year}
			}
		} else {
			// If search is empty, append everything
			gallery.Years[year] = Index{[]Post{}, append(gallery.Years[year].Pictures, pic), year}
		}
	}

	if en {
		template.Must(template.ParseFiles("web/en/gallery.html")).Execute(w, gallery)
	} else {
		template.Must(template.ParseFiles("web/gallery.html")).Execute(w, gallery)
	}
}

