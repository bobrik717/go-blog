package main

import (
	"fmt"
	"net/http"
	"html/template"
	"go-book/blog/models"
)

var posts map[string]*models.Post

func main() {
	fmt.Println("Listnening on port :3000")

	posts = make(map[string]*models.Post,0)

	http.Handle("/assets/", http.StripPrefix("/assets/",http.FileServer(http.Dir("./assets/"))))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/write", writeHandler)
	http.HandleFunc("/savePost", savePostHandler)

	http.ListenAndServe(":3000",nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("views/index.html","views/header.html","views/footer.html")
	if err != nil {
		fmt.Fprintf(w,err.Error())
	}
	temp.ExecuteTemplate(w,"index",posts)
}

func writeHandler(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("views/write.html","views/header.html","views/footer.html")
	if err != nil {
		fmt.Fprintf(w,err.Error())
	}
	temp.ExecuteTemplate(w,"write",nil)
}

func savePostHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	content := r.FormValue("content")

	post := models.NewPost(id,title,content)
	posts[post.Id] = post

	http.Redirect(w,r,"/",302)
}