package models

type Library struct {
	ID      string `json:"id" bson:"_id"`
	Name    string `json:"name" bson:"name"`
	Address string `json:"address" bson:"address"`
	Year    string `json:"year" bson:"year"`
	Books  []Book `json:"books" bson:"books"`
}
