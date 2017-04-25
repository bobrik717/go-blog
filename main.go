package main

import (
	"net/http"
	"./models"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/russross/blackfriday"
	"html/template"
)

var posts map[string]*models.Post

func main() {
	posts = make(map[string]*models.Post,0)

	m := models.MyMartiniClassic{*martini.Classic()}

	unescapeFuncMap := template.FuncMap{"unescape": models.Unescape}

	m.Use(render.Renderer(render.Options{
		Directory: "views", // Specify what path to load the templates from.
		Layout: "layout", // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates.
		Funcs: []template.FuncMap{unescapeFuncMap}, // Specify helper function maps for templates to access.
		Charset: "UTF-8", // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON: true, // Output human readable JSON
	}))

	options := martini.StaticOptions{Prefix:"assets"}
	m.Use(martini.Static("assets",options))

	m.Get("/", homeHandler)
	m.Get("/write", writeHandler)
	m.Get("/edit/:id", editHandler)
	m.Post("/savePost", savePostHandler)
	m.Get("/deletePost/:id", deletePostHandler)
	m.Post("/gethtml", getHtmlHandler)

	m.RunCustom("3001")
}

func homeHandler(rnd render.Render) {
	rnd.HTML(200, "index", posts)
}

func writeHandler(rnd render.Render) {
	rnd.HTML (200,"write",nil)
}

func editHandler(rnd render.Render, params martini.Params) {
	id := params["id"]

	pots, isFound := posts[id]
	if !isFound {
		rnd.Error(404)
	}

	rnd.HTML(200,"write", pots)
}

func savePostHandler(rnd render.Render, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	contentMarkDown := r.FormValue("content")
	contentHtml  := string(blackfriday.MarkdownBasic([]byte(contentMarkDown)))

	if id == "" {
		id = models.GenerateId()
		post := models.NewPost(id, title, contentHtml, contentMarkDown)
		posts[post.Id] = post
	} else {
		_, isFound := posts[id]
		if !isFound {
			rnd.Error(404)
		}

		posts[id].Title = title
		posts[id].ContentHtml = contentHtml
		posts[id].ContentMarkDown = contentMarkDown
	}
	rnd.Redirect("/",302)
}

func deletePostHandler(rnd render.Render, params martini.Params) {
	id := params["id"]
	_, isFound := posts[id]
	if !isFound {
		rnd.Error(404)
	}
	delete(posts, id)
	rnd.Redirect("/",302)
}

func getHtmlHandler(rnd render.Render, r *http.Request) {
	md := r.FormValue("md")
	htmlBytes := blackfriday.MarkdownBasic([]byte(md))

	rnd.JSON(200, map[string]interface{} {"html": string(htmlBytes)})
}