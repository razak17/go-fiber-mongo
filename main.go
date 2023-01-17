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

	// start server
	port := os.Getenv("PORT")

	app.Listen(":" + port)

	return nil
}

func loadEnv() error {
	goEnv := os.Getenv("GO_ENV")
	if goEnv == "" {
		err := godotenv.Load()
		if err != nil {
			return err
		}
	}
	return nil
}
