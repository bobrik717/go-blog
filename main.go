package main

import (
	"fmt"
	"net/http"
	"html/template"
	"./models"
)

var posts map[string]*models.Post

func main() {
	fmt.Println("Listnening on port :3000")

	posts = make(map[string]*models.Post,0)

	http.Handle("/assets/", http.StripPrefix("/assets/",http.FileServer(http.Dir("./assets/"))))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/write", writeHandler)
	http.HandleFunc("/edit", editHandler)
	http.HandleFunc("/savePost", savePostHandler)
	http.HandleFunc("/deletePost", deletePostHandler)

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

func editHandler(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("views/write.html","views/header.html","views/footer.html")
	if err != nil {
		fmt.Fprintf(w,err.Error())
	}

	id := r.FormValue("id")

	pots, isFound := posts[id]
	if !isFound {
		http.NotFound(w,r)
	}

	temp.ExecuteTemplate(w,"write", pots)
}

func savePostHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	content := r.FormValue("content")

	if id == "" {
		id = models.GenerateId()
		post := models.NewPost(id,title,content)
		posts[post.Id] = post
	} else {
		_, isFound := posts[id]
		if !isFound {
			http.NotFound(w,r)
		}

		posts[id].Title = title
		posts[id].Content = content
	}
	http.Redirect(w,r,"/",302)
}

func deletePostHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	_, isFound := posts[id]
	if !isFound {
		http.NotFound(w,r)
	}
	delete(posts, id)
	http.Redirect(w,r,"/",302)
}