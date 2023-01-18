package handlers

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/razak17/go-fiber-mongo/database"
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

	// get the library id
	objectId, err := primitive.ObjectIDFromHex(newBook.LibraryId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid libraryId"})
	}

	// update the library
	filter := bson.D{{Key: "_id", Value: objectId}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "books", Value: result.InsertedID}}}}

	collection := database.GetDBCollection("libraries")
	_, libraryErr := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get library",
			"message": libraryErr.Error(),
		})
	}

	// return the library
	return c.Status(fiber.StatusCreated).SendString(fmt.Sprintf("Created book with id: %s", result.InsertedID))
}
