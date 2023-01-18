package models

type Book struct {
	ID        string `json:"id" bson:"_id"`
	LibraryId string `json:"libraryId" bson:"libraryId"`
	Title     string `json:"title" bson:"title"`
	Author    string `json:"author" bson:"author"`
	ISBN      string `json:"isbn" bson:"isbn"`
}
