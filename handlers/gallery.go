package handlers

import (
	"net/http"
	"os"
	"html/template"
	"strings"
)


func Pictures(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "content/pictures/" + r.PathValue("file"))
}


func Pics(w http.ResponseWriter, r *http.Request) {
	image := r.PathValue("pic")

	en := (r.URL.Path == "/en/pics/" + image)
		
	// Checks if the file exists
	_, err := os.Stat("content/pictures/" + image)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	tm := GetTime(("content/pictures/" + image))
	title := image[:len(image)-4]

	if en {
		template.Must(template.ParseFiles("web/en/picture.html")).Execute(w, Picture{title,"", tm.Format("2/1/2006 - 15:04"), image, tm, GetHtml(("content/pictures/en/" + image + ".md"))})
	} else {
		template.Must(template.ParseFiles("web/picture.html")).Execute(w, Picture{title,"", tm.Format("2/1/2006 - 15:04"), image, tm, GetHtml(("content/pictures/" + image + ".md"))})
	}
}


func Gallery(w http.ResponseWriter, r *http.Request) {
	en := (r.URL.Path == "/en/pics/")

	search := r.URL.Query().Get("q")


	var gallery = struct{Title string; Years map[string]Index}{"Galeria", make(map[string]Index)}

	if search != "" {
		if en {
			gallery.Title = "Results for \"" + search + "\""
		} else {
			gallery.Title = "Resultados para \"" + search + "\""
		}
	}

	var list []Picture

	if en {
		list = GetPictures("content/pictures", "/en")
	} else {
		list = GetPictures("content/pictures", "")
	}

	for _, pic := range list {
		year := pic.Time.Format("2006")

		if search != "" {
			search := strings.ToLower(search)

			// If search is not empty, append only the matching content
			if strings.Contains(strings.ToLower(pic.Title), search) || strings.Contains(strings.ToLower(pic.Summary), search) || strings.Contains(pic.Time.Format("2006"), search) {
				gallery.Years[year] = Index{[]Page{}, append(gallery.Years[year].Pictures, pic), year}
			}
		} else {
			// If search is empty, append everything
			gallery.Years[year] = Index{[]Page{}, append(gallery.Years[year].Pictures, pic), year}
		}
	}

	if en {
		template.Must(template.ParseFiles("web/en/gallery.html")).Execute(w, gallery)
	} else {
		template.Must(template.ParseFiles("web/gallery.html")).Execute(w, gallery)
	}
}

