package handlers

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/razak17/go-fiber-mongo/database"
)

type libraryDTO struct {
	Name    string `json:"name" bson:"name"`
	Address string `json:"address" bson:"address"`
}

func CreateLibraryHandler(c *fiber.Ctx) error {
	// validate the body
	newLibrary := new(libraryDTO)
	if err := c.BodyParser(newLibrary); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid body",
		})
	}

	fmt.Println(newLibrary)

	// create the library
	collection := database.GetDBCollection("libraries")
	result, err := collection.InsertOne(context.TODO(), newLibrary)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create book",
			"message": err.Error(),
		})
	}

	// return the book
	return c.Status(201).JSON(fiber.Map{
		"result": result.InsertedID,
	})
}
