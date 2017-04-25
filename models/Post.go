package models

type Post struct {
	Id string
	Title string
	ContentHtml string
	ContentMarkDown string
}

func NewPost(id, title, contentHtml, contentMarkDown string) *Post {
	return &Post{id,title,contentHtml, contentMarkDown}
}