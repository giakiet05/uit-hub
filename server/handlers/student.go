package handlers

import (
	"server/data"

	"github.com/gofiber/fiber/v2"
)

// GET /students
func GetStudents(c *fiber.Ctx) error {
	return c.JSON(data.Students)
}

// GET /students/:id
func GetStudentByID(c *fiber.Ctx) error {
	id := c.Params("id")

	for _, student := range data.Students {
		if student.ID == id {
			return c.JSON(student)
		}
	}

	return c.Status(404).JSON(fiber.Map{
		"error": "Student not found",
	})
}
