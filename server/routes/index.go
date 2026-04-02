package routes

import (
	"server/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Get("/students", handlers.GetStudents)
	api.Get("/students/:id", handlers.GetStudentByID)
}
