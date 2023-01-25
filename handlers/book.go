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

	// push book into the library
	filter := bson.D{{Key: "_id", Value: objectId}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "books", Value: result.InsertedID}}}}
	collection := database.GetDBCollection("libraries")
	_, libraryErr := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to push book",
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

type updateBookDTO struct {
	LibraryId string `json:"libraryId,omitempty" bson:"libraryId,omitempty"`
	Title     string `json:"title,omitempty" bson:"title,omitempty"`
	Author    string `json:"author,omitempty" bson:"author,omitempty"`
	ISBN      string `json:"isbn,omitempty" bson:"isbn,omitempty"`
}

func UpdateBook(c *fiber.Ctx) error {
	b := new(updateBookDTO)
	if err := c.BodyParser(b); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	// get the id
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id is required"})
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	collection := database.GetDBCollection("books")
	_, err = collection.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: objectId}}, bson.D{{Key: "$set", Value: b}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update book",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).SendString(fmt.Sprintf("Updated book with id: %s", id))
}

func DeleteBook(c *fiber.Ctx) error {
	// get the id
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id is required"})
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	// Get the book
	book := models.Book{}
	bookCollection := database.GetDBCollection("books")
	err = bookCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objectId}}).Decode(&book)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get book",
			"message": err.Error(),
		})
	}

	libraryId, err := primitive.ObjectIDFromHex(book.LibraryId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	// remove book from the library
	collection := database.GetDBCollection("libraries")
	filter := bson.D{{Key: "_id", Value: libraryId}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "books", Value: objectId}}}}
	_, libraryErr := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to remove book",
			"message": libraryErr.Error(),
		})
	}

	// delete the book
	filter = bson.D{{Key: "_id", Value: objectId}}
	_, err = bookCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to delete book",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).SendString(fmt.Sprintf("Deleted book with id: %s", libraryId))
}
