package actions

import (
	"go-book/blog/models"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"html/template"
	"net/http"
	"go-book/blog/session"
)

func App() {
	app := models.MyMartiniClassic{*martini.Classic()}

	unescapeFuncMap := template.FuncMap{"unescape": models.Unescape}

	app.Use(render.Renderer(render.Options{
		Directory: "views", // Specify what path to load the templates from.
		Layout: "layout", // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates.
		Funcs: []template.FuncMap{unescapeFuncMap}, // Specify helper function maps for templates to access.
		Charset: "UTF-8", // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON: true, // Output human readable JSON
	}))

	options := martini.StaticOptions{Prefix:"assets"}
	app.Use(martini.Static("assets",options))

	app.Use(func(r *http.Request, rnd render.Render) {
		_, err := r.Cookie(session.COOKIE_NAME)
		if err != nil && r.RequestURI != "/login" {
			rnd.Redirect("/login",302)
		}
	})

	app.Get("/", IndexHandler)
	app.Get("/login", LoginIndexHandler)
	app.Post("/login", LoginHandler)
	app.Get("/write", WriteHandler)
	app.Get("/edit/:id", EditHandler)
	app.Post("/savePost", SavePostHandler)
	app.Get("/deletePost/:id", DeletePostHandler)
	app.Post("/gethtml", GetHtmlHandler)

	app.RunCustom("3001")
}