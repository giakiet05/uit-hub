package controller

import (
	"net/http"

	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/apperror"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/dto"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/middleware"
	"github.com/gin-gonic/gin"
)

func (c *UitController) GetDeadlines(ctx *gin.Context) {

	studentID, ok := ctx.Get(middleware.StudentIDKey)
	if !ok {
		dto.SendError(ctx, apperror.ErrUnauthorized)
		return
	}

	res, appErr := c.service.GetDeadlines(studentID.(string))
	if appErr != nil {
		dto.SendError(ctx, *appErr)
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, res.Data)
}

func (c *UitController) GetMaterials(ctx *gin.Context) {

	studentID, ok := ctx.Get(middleware.StudentIDKey)
	if !ok {
		dto.SendError(ctx, apperror.ErrUnauthorized)
		return
	}

	courseID := ctx.Param("courseId")

	res, appErr := c.service.GetMaterials(studentID.(string), courseID)
	if appErr != nil {
		dto.SendError(ctx, *appErr)
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, res.Data)
}

func (c *UitController) GetAssignments(ctx *gin.Context) {

	studentID, ok := ctx.Get(middleware.StudentIDKey)
	if !ok {
		dto.SendError(ctx, apperror.ErrUnauthorized)
		return
	}

	courseID := ctx.Param("courseId")

	res, appErr := c.service.GetAssignments(studentID.(string), courseID)
	if appErr != nil {
		dto.SendError(ctx, *appErr)
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, res.Data)
}

func (c *UitController) SubmitAssignment(ctx *gin.Context) {

	studentID, ok := ctx.Get(middleware.StudentIDKey)
	if !ok {
		dto.SendError(ctx, apperror.ErrUnauthorized)
		return
	}

	assignmentID := ctx.Param("assignmentId")

	var req dto.AssignmentSubmissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, apperror.ErrBadRequest)
		return
	}

	res, appErr := c.service.SubmitAssignment(studentID.(string), assignmentID, req)
	if appErr != nil {
		dto.SendError(ctx, *appErr)
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, res.Data)
}
