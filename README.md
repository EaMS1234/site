# site

My personal website featuring my blog, gallery and portfolio!
Databaseless web app that reads markdown and image files directly from the host's filesystem, everything running in Go.

## Usage
### Building and Running
- Install the dependencies with `go mod tidy`
- Build with `go build .`
- Run with `./site`

### Docker
You can deploy the server as a docker container by using the provided dockerfile:
- Build the image with `docker build -t site .`
- Check the container image with `docker run -p 8080:8080 site` and open it in your browser on `localhost:8080/`

Please note that in orther to add content, you must have a way of accesing the container's `/site/content` folder.
You can try adding `-v <folder>:/site/content` to the `docker run` command, where `<folder>` is the path of a folder acessible by the host. 

### Adding content
The server reads the `content` folder looking for `.md`, `.png` and `.jpg` files and parses/serves them in nice-loooking pages provideded by the HTML templates on `web`.
It can also serve static assets and whole files contained inside the `static` folder.

The basic structure of the `content` directory is:
- `content/about.md`: The "about" page for the main language (portuguese).
- `content/about.en.md`: The "about" page for the english language.
- `content/posts/`: Location of posts, as `.md` files, for the main language. The name of the file itself is the title of the post.
- `content/posts/en/`: Location of posts for the english language.
- `content/pictures/`: Location of pictures, as `.png` and `.jpg` files. Captions can be providade in markdown by sharing its name with the picture: `example.png.md`.
- `content/pictures/en/`: Location of pictures' captions as `.md` files for the english language.
- `content/static/`: Serves all files within it as-is directly trough HTTP.

## Roadmap
- [x] Read and parse content from disk
- [x] Blog
- [x] Gallery
- [x] Multilingual support (portuguese & english)
- [ ] Accessibility features
- [ ] Portfolio page for projects
- [ ] Support for the russian language
- [ ] Support for tags, categories and picture folders
- [ ] RSS feed
