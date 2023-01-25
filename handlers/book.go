package handlers

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/razak17/go-fiber-mongo/database"
	"github.com/razak17/go-fiber-mongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type bookDTO struct {
	LibraryId string `json:"libraryId" bson:"libraryId"`
	Title     string `json:"title" bson:"title"`
	Author    string `json:"author" bson:"author"`
	ISBN      string `json:"isbn" bson:"isbn"`
}

func CreateBook(c *fiber.Ctx) error {
	// validate the body
	newBook := new(bookDTO)
	if err := c.BodyParser(newBook); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	//  create book doc
	result, err := database.GetDBCollection("books").InsertOne(context.TODO(), newBook)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create book"})
	}

	// get the book id
	objectId, err := primitive.ObjectIDFromHex(newBook.LibraryId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid libraryId"})
	}

	// update the book
	filter := bson.D{{Key: "_id", Value: objectId}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "books", Value: result.InsertedID}}}}

	collection := database.GetDBCollection("libraries")
	_, libraryErr := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get book",
			"message": libraryErr.Error(),
		})
	}

	// return the book
	return c.Status(fiber.StatusCreated).SendString(fmt.Sprintf("Created book with id: %s", result))
}

func GetBooks(c *fiber.Ctx) error {
	collection := database.GetDBCollection("books")

	books := make([]models.Book, 0)
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get books",
			"message": err.Error(),
		})
	}

	// iterate over the cursor
	for cursor.Next(context.TODO()) {
		book := models.Book{}
		err := cursor.Decode(&book)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		books = append(books, book)
	}

	defer cursor.Close(context.TODO())
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": books})
}

func GetBook(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id is required"})
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	book := models.Book{}
	collection := database.GetDBCollection("books")
	err = collection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objectId}}).Decode(&book)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get book",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": book})
}

func UpdateBook(c *fiber.Ctx) error {
	return c.SendString("Update book")
}

func DeleteBook(c *fiber.Ctx) error {
	return c.SendString("Delete book")
}
