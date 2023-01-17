package handlers

import "github.com/gofiber/fiber/v2"

func HealthCheckHandler(c *fiber.Ctx) error {
	return c.SendString("OK")
}

func TestHandler(c *fiber.Ctx) error {
	return c.SendString("Hello World")
}

// app.Post("/hello", func(c *fiber.Ctx) error {
// 	doc := bson.M{"Atonement": "Ian McEwan"}
// 	collection := database.GetDBCollection("books")
// 	result, err := collection.InsertOne(context.TODO(), doc)
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}
//
// 	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
//
// 	// return the book
// 	return c.Status(201).JSON(fiber.Map{
// 		"result": result,
// 	})
// })
