package controller

import (
	"net/http"

	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/apperror"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/dto"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/middleware"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/service"
	"github.com/gin-gonic/gin"
)

type StudentController struct {
	service service.StudentService
}

func NewStudentController(service service.StudentService) *StudentController {
	return &StudentController{service: service}
}

func (c *StudentController) GetStudents(ctx *gin.Context) {
	students := c.service.GetStudents()
	dto.SendSuccess(ctx, http.StatusOK, dto.NewStudentResponses(students))
}

func (c *StudentController) GetProfile(ctx *gin.Context) {
	studentID, ok := ctx.Get(middleware.StudentIDKey)
	if !ok {
		dto.SendError(ctx, apperror.ErrUnauthorized)
		return
	}

	student, found := c.service.GetStudentByID(studentID.(string))
	if !found {
		dto.SendError(ctx, apperror.ErrNotFound)
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, dto.NewStudentResponse(*student))
}

func (c *StudentController) GetStudentByID(ctx *gin.Context) {
	student, found := c.service.GetStudentByID(ctx.Param("id"))
	if !found {
		err := apperror.ErrNotFound
		err.Message = "Student not found"
		dto.SendError(ctx, err)
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, dto.NewStudentResponse(*student))
}
