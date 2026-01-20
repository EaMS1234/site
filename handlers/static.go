package handlers

import "net/http"


// Routes static assets like stylesheets and scripts
func InitAssets(content_dir string) {
	styles := http.FileServer(http.Dir("web/styles/"))
	// scripts := http.FileServer(http.Dir("web/scripts/"))
	assets := http.FileServer(http.Dir("web/assets/"))

	http.Handle("/styles/", http.StripPrefix("/styles/", styles))
	// http.Handle("/scripts/", http.StripPrefix("/scripts/", scripts))
	http.Handle("/assets/", http.StripPrefix("/assets/", assets))

	// Content
	pictures := http.FileServer(http.Dir(content_dir + "/pictures"))
	static := http.FileServer(http.Dir(content_dir + "/static"))
	http.Handle("/pictures/", http.StripPrefix("/pictures/", pictures))
	http.Handle("/static/", http.StripPrefix("/static/", static))
}

