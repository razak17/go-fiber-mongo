package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/razak17/go-fiber-mongo/database"
	"github.com/razak17/go-fiber-mongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type libraryDTO struct {
	Name    string   `json:"name" bson:"name"`
	Address string   `json:"address" bson:"address"`
	Empty   []string `json:"no_exists" bson:"books"`
}

func CreateLibraryHandler(c *fiber.Ctx) error {
	// validate the body
	newLibrary := new(libraryDTO)
	if err := c.BodyParser(newLibrary); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid body",
		})
	}

	newLibrary.Empty = make([]string, 0)

	// create the library
	collection := database.GetDBCollection("libraries")
	result, err := collection.InsertOne(context.TODO(), newLibrary)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create library",
			"message": err.Error(),
		})
	}

	// return the library
	return c.Status(201).JSON(fiber.Map{
		"result": result.InsertedID,
	})
}

func GetLibraries(c *fiber.Ctx) error {
	collection := database.GetDBCollection("libraries")

	// find all libraries
	libraries := make([]models.Library, 0)
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to get libraries",
			"message": err.Error(),
		})
	}

	// iterate over the cursor
	for cursor.Next(c.Context()) {
		library := models.Library{}
		err := cursor.Decode(&library)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		libraries = append(libraries, library)
	}

	return c.Status(200).JSON(fiber.Map{"data": libraries})
}

func GetLibrary(c *fiber.Ctx) error {
	collection := database.GetDBCollection("libraries")

	// find the library
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "id is required",
		})
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	library := models.Library{}

	err = collection.FindOne(c.Context(), bson.M{"_id": objectId}).Decode(&library)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{"data": library})
}

type updateLibraryDTO struct {
	Name    string `json:"name" bson:"name"`
	Address string `json:"address" bson:"address"`
}

func UpdateLibrary(c *fiber.Ctx) error {
	// validate the body
	l := new(updateLibraryDTO)
	if err := c.BodyParser(l); err != nil {
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

	// update the book
	filter := bson.D{{Key: "_id", Value: objectId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: l.Name}, {Key: "address", Value: l.Address}}}}

	result, err := database.GetDBCollection("libraries").UpdateOne(c.Context(), filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update book",
			"message": err.Error(),
		})
	}

	// return the book
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"result": result})
}
