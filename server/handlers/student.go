package handlers

import "github.com/gofiber/fiber/v2"

// GET /api/students
func (h *Handler) GetStudents(c *fiber.Ctx) error {
	resp, appErr := h.uc.ListStudents()
	return h.respond(c, resp, appErr)
}

// GET /api/students/:id
func (h *Handler) GetStudentByID(c *fiber.Ctx) error {
	id := c.Params("id")
	resp, appErr := h.uc.GetStudentByID(id)
	return h.respond(c, resp, appErr)
}
