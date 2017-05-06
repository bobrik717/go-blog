package main

import (
	"go-book/blog/actions"
	"gopkg.in/mgo.v2"
	"go-book/blog/session"
)

func main() {
	dbConnect, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer dbConnect.Close()

	session.SessionMain = session.NewSession()

	actions.PostsCollection = dbConnect.DB("blog").C("posts")

	actions.App()
}