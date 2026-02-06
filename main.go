package main

import (
	"log"
	"net/http"
	"os"

	"github.com/eams1234/site/handlers"
)


func main() {
	content_dir := "content/"

	// Checks if content contains the necessary folders
	_, err1 := os.Stat(content_dir + "/posts")
	_, err2 := os.Stat(content_dir + "/pictures")
	_, err3 := os.Stat(content_dir + "/static")

	if err1 != nil {
		os.Mkdir((content_dir + "/posts"), 0755)
		os.Mkdir((content_dir + "/posts/en"), 0755)
	}

	if err2 != nil {
		os.Mkdir((content_dir + "/pictures"), 0755)
		os.Mkdir((content_dir + "/pictures/en"), 0755)
	}

	if err3 != nil {
		os.Mkdir((content_dir + "/static"), 0755)
	}

	mux := http.NewServeMux()

	// Loads the content
	handlers.GetPosts()
	handlers.GetPictures()

	// Static files
	mux.HandleFunc("/assets/{file}", handlers.Assets)
	mux.HandleFunc("/styles/{file}", handlers.Styles)
	mux.HandleFunc("/scripts/{file}", handlers.Scripts)
	mux.HandleFunc("/static/{file}", handlers.Static)
	mux.HandleFunc("/img/{file}", handlers.Images)

	// Index
	mux.HandleFunc("/", handlers.InitIndex)

	// About
	mux.HandleFunc("/sobre/", handlers.About)
	mux.HandleFunc("/en/about/", handlers.About)

	// Blog & gallery
	mux.HandleFunc("/artigos/", handlers.Blog)
	mux.HandleFunc("/en/blog/", handlers.Blog)
	mux.HandleFunc("/galeria/", handlers.Gallery)
	mux.HandleFunc("/en/pictures/", handlers.Gallery)

	// Individual posts
	mux.HandleFunc("/artigos/{post}", handlers.Posts)
	mux.HandleFunc("/en/blog/{post}", handlers.Posts)
	mux.HandleFunc("/galeria/{pic}", handlers.Pictures)
	mux.HandleFunc("/en/pictures/{pic}", handlers.Pictures)

	log.Output(1, "Serving on port 8080")

	http.ListenAndServe(":8080", mux)
}

