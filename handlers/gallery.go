package handlers

import (
	"net/http"
	"os"
	"html/template"
	"strings"
)


func InitGallery(content_dir string) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		en := (r.URL.Path == "/en/pics/")

		image := r.URL.Query().Get("i")
		search := r.URL.Query().Get("q")

		if image != "" {
			// This means it's not a valid image, as the .jpg and .png endings are
			// both exactly 4 characters long.
			if len(image) <= 4 {
				if en {
					handle404(w, r, "en")
					return
				} else {
					handle404(w, r, "")
					return
				}
			}

			title := image[:len(image)-4]
		
			// Checks if the file exists
			_, err := os.Stat(content_dir + "/pictures/" + image)
			if err != nil {
				if en {
					handle404(w, r, "en")
				} else {
					handle404(w, r, "")
				}
				return
			}

			tm := GetTime((content_dir + "/pictures/" + image))

			var desc template.HTML

			if en {
				desc = GetHtml((content_dir + "/pictures/en/" + image + ".md"))
				template.Must(template.ParseFiles("web/en/picture.html")).Execute(w, Picture{title,"", tm.Format("2/1/2006 - 15:04"), image, tm, desc})
			} else {
				desc = GetHtml((content_dir + "/pictures/" + image + ".md"))
				template.Must(template.ParseFiles("web/picture.html")).Execute(w, Picture{title,"", tm.Format("2/1/2006 - 15:04"), image, tm, desc})
			}
		} else {
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
				list = GetPictures(content_dir + "/pictures", "/en")
			} else {
				list = GetPictures(content_dir + "/pictures", "")
			}

			for _, pic := range list {
				year := pic.Time.Format("2006")
				
				if search != "" {

					// If search is not empty, append only the matching content
					if strings.Contains(pic.Title, search) || strings.Contains(pic.Summary, search) || strings.Contains(pic.Time.Format("2006"), search) {
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
	}

	http.HandleFunc("/galeria/", handler)
	http.HandleFunc("/en/pics/", handler)
}

