package actions

import (
	"go-book/blog/db/documents"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"net/http"
	"fmt"
	"go-book/blog/models"
	"gopkg.in/mgo.v2"
	"go-book/blog/session"
)

var PostsCollection *mgo.Collection

func IndexHandler(rnd render.Render, r *http.Request) {
	postDocuments := []documents.PostDocument{}
	PostsCollection.Find(nil).All(&postDocuments)

	id, err := r.Cookie(session.COOKIE_NAME)
	if err == nil {
		fmt.Println(session.SessionMain.Get(id.Value))
	}

	posts := []*models.Post{}
	for _,doc := range postDocuments {
		post := models.NewPost(doc.Id,doc.Title,doc.ContentHtml,doc.ContentMarkDown)
		posts = append(posts, post)
	}
	rnd.HTML(200, "index", posts)
}

func WriteHandler(rnd render.Render) {
	rnd.HTML (200,"write",nil)
}

func EditHandler(rnd render.Render, params martini.Params) {
	id := params["id"]
	doc := documents.PostDocument{}
	err := PostsCollection.FindId(id).One(&doc)
	if err != nil {
		rnd.Error(404)
	}
	post := models.NewPost(doc.Id,doc.Title,doc.ContentHtml,doc.ContentMarkDown)
	rnd.HTML(200,"write", post)
}

func SavePostHandler(rnd render.Render, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	contentMarkDown := r.FormValue("content")
	contentHtml  := models.ConvertMarkDownToHtml(contentMarkDown)

	postDocument := documents.PostDocument{id,title,contentHtml,contentMarkDown}

	if id == "" {
		id = models.GenerateId()
		postDocument.Id = id
		PostsCollection.Insert(postDocument)
	} else {
		PostsCollection.UpdateId(id, postDocument)
	}
	rnd.Redirect("/",302)
}

func DeletePostHandler(rnd render.Render, params martini.Params) {
	id := params["id"]
	PostsCollection.RemoveId(id)
	rnd.Redirect("/",302)
}

func GetHtmlHandler(rnd render.Render, r *http.Request) {
	md := r.FormValue("md")
	rnd.JSON(200, map[string]interface{} {"html": models.ConvertMarkDownToHtml(md)})
}