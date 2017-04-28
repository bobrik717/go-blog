package main

import (
	"net/http"
	"./models"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"html/template"
	"gopkg.in/mgo.v2"
	"./db/documents"
	"./session"
	"time"
	"fmt"
)

var postsCollection *mgo.Collection
var sessionMain *session.Session

const (
	COOKIE_NAME = "sessionId"
)

func main() {
	dbConnect, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer dbConnect.Close()

	sessionMain = session.NewSession()

	postsCollection = dbConnect.DB("blog").C("posts")

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
	m.Get("/login", loginIndexHandler)
	m.Post("/login", loginHandler)
	m.Get("/write", writeHandler)
	m.Get("/edit/:id", editHandler)
	m.Post("/savePost", savePostHandler)
	m.Get("/deletePost/:id", deletePostHandler)
	m.Post("/gethtml", getHtmlHandler)

	m.RunCustom("3001")
}

func loginIndexHandler(rnd render.Render) {
	rnd.HTML (200,"login",nil)
}

func loginHandler(rnd render.Render, r *http.Request, w http.ResponseWriter) {
	login := r.FormValue("login")
	//password := r.FormValue("password")

	sessionId := sessionMain.Init(login)

	cookie := &http.Cookie{
		Name: COOKIE_NAME,
		Value: sessionId,
		Expires: time.Now().Add( 5 * time.Minute),
	}

	http.SetCookie(w, cookie)

	rnd.Redirect("/")
}

func homeHandler(rnd render.Render, r *http.Request) {
	postDocuments := []documents.PostDocument{}
	postsCollection.Find(nil).All(&postDocuments)

	id, err := r.Cookie(COOKIE_NAME)
	if err == nil {
		fmt.Println(sessionMain.Get(id.Value))
	}

	posts := []*models.Post{}
	for _,doc := range postDocuments {
		post := models.NewPost(doc.Id,doc.Title,doc.ContentHtml,doc.ContentMarkDown)
		posts = append(posts, post)
	}
	rnd.HTML(200, "index", posts)
}

func writeHandler(rnd render.Render) {
	rnd.HTML (200,"write",nil)
}

func editHandler(rnd render.Render, params martini.Params) {
	id := params["id"]
	doc := documents.PostDocument{}
	err := postsCollection.FindId(id).One(&doc)
	if err != nil {
		rnd.Error(404)
	}
	post := models.NewPost(doc.Id,doc.Title,doc.ContentHtml,doc.ContentMarkDown)
	rnd.HTML(200,"write", post)
}

func savePostHandler(rnd render.Render, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	contentMarkDown := r.FormValue("content")
	contentHtml  := models.ConvertMarkDownToHtml(contentMarkDown)

	postDocument := documents.PostDocument{id,title,contentHtml,contentMarkDown}

	if id == "" {
		id = models.GenerateId()
		postDocument.Id = id
		postsCollection.Insert(postDocument)
	} else {
		postsCollection.UpdateId(id, postDocument)
	}
	rnd.Redirect("/",302)
}

func deletePostHandler(rnd render.Render, params martini.Params) {
	id := params["id"]
	postsCollection.RemoveId(id)
	rnd.Redirect("/",302)
}

func getHtmlHandler(rnd render.Render, r *http.Request) {
	md := r.FormValue("md")
	rnd.JSON(200, map[string]interface{} {"html": models.ConvertMarkDownToHtml(md)})
}