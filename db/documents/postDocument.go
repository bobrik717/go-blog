package documents

type PostDocument struct {
	Id 		string `bson:"_id,omitempty"`
	Title 		string
	ContentHtml 	string
	ContentMarkDown string
}