package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"taskManager/db"
	"taskManager/logger"
	"taskManager/manager"
)

const _URI = "mongodb://localhost:27017"

func hello(c *fiber.Ctx) error {
	return c.SendString("Hello World")
}

func main() {
	db.Connect(_URI)
	logger.InitLogger()

	defer db.Disconnect()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	app.Get("/", hello) // for testing
	app.Get("/tasks", manager.ViewTasks)
	app.Post("/task", manager.AddTask)
	app.Put("/task/:id", manager.CompleteTask)
	app.Delete("/task/:id", manager.DeleteTask)
	app.Put("/filter/:filter", manager.UseFilter)
	app.Get("/filter", manager.GetFilter)

	err := app.Listen(":3131")
	if err != nil {
		fmt.Println("Die")
	}
}
