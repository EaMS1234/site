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


	handlers.InitAssets(content_dir)
	handlers.InitIndex(content_dir)
	handlers.InitAbout(content_dir)
	handlers.InitPosts(content_dir)
	handlers.InitGallery(content_dir)

	log.Output(1, "Serving on port 8080")

	http.ListenAndServe(":8080", nil)
}

