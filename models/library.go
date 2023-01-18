package models

type Library struct {
	ID      string   `json:"id" bson:"_id"`
	Name    string   `json:"name" bson:"name"`
	Address string   `json:"address" bson:"address"`
	Books   []string `json:"books" bson:"books"`
}
