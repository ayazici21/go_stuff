package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"taskManager/auth"
	"taskManager/db"
	"taskManager/logger"
	"taskManager/manager"
	"taskManager/user"
)

const _URI = "mongodb://localhost:27017"

func routes(app *fiber.App) {
	app.Get("/tasks/*", manager.ViewTasks)
	app.Post("/task", manager.AddTask)
	app.Put("/task/:id", manager.CompleteTask)
	app.Delete("/task/:id", manager.DeleteTask)

	app.Post("/user/register", user.Register)
	app.Post("/user/login", user.Login)
	app.Post("/user/logout", user.Logout)

}

func main() {
	db.Connect(_URI)
	log := logger.New()

	defer db.Disconnect()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	app.Use("/task/*", auth.ExtractToken, auth.StaleCheck)
	app.Use("/tasks", auth.ExtractToken, auth.StaleCheck)
	app.Use("/user/logout", auth.ExtractToken, auth.StaleCheck)

	routes(app)

	err := app.Listen(":3131")

	if err != nil {
		log.Error("I dieded")
		log.Error("%s", err)
		panic(err)
	}
}
