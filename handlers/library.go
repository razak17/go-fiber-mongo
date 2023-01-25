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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	newLibrary.Empty = make([]string, 0)

	// create the library
	collection := database.GetDBCollection("libraries")
	result, err := collection.InsertOne(context.TODO(), newLibrary)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create library",
			"message": err.Error(),
		})
	}

	// return the library
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"result": result.InsertedID})
}

func GetLibraries(c *fiber.Ctx) error {
	collection := database.GetDBCollection("libraries")

	// find all libraries
	libraries := make([]models.Library, 0)
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get libraries",
			"message": err.Error(),
		})
	}

	// iterate over the cursor
	for cursor.Next(context.TODO()) {
		library := models.Library{}
		err := cursor.Decode(&library)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		libraries = append(libraries, library)
	}

	defer cursor.Close(context.TODO())
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": libraries})
}

func GetLibrary(c *fiber.Ctx) error {
	// find the library
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id is required"})
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	library := models.Library{}

	collection := database.GetDBCollection("libraries")
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(&library)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": library})
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

	// update the library
	filter := bson.D{{Key: "_id", Value: objectId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: l.Name}, {Key: "address", Value: l.Address}}}}

	collection := database.GetDBCollection("libraries")
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update library",
			"message": err.Error(),
		})
	}

	// return the library
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"result": result})
}

func DeleteLibrary(c *fiber.Ctx) error {
	// get the id
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id is required"})
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	// delete the library
	collection := database.GetDBCollection("libraries")
	filter := bson.D{{Key: "_id", Value: objectId}}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to delete library",
			"message": err.Error(),
		})
	}

	// delete books belonging to the library
	_, err = database.GetDBCollection("books").DeleteMany(context.TODO(), bson.D{{Key: "libraryId", Value: id}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to delete library books",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"result": result})
}

func GetLibraryBooks(c *fiber.Ctx) error {
	return c.SendString("GetLibraryBooks")
}
