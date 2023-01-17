package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/razak17/go-fiber-mongo/database"
)

func main() {
	err := initApp()
	if err != nil {
		panic(err)
	}
}

func initApp() error {
	// load env
	err := loadEnv()
	if err != nil {
		return err
	}

	// setup mongoDB
	err = database.ConnectDB()
	if err != nil {
		return err
	}

	// defer closing db
	defer database.CloseDB()

	// create app
	app := generateApp()

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

	// start server
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "8080"
	}
	app.Listen(":" + port)

	return nil
}

func loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}
