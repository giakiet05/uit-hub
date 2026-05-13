package controller

import (
	"net/http"

	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/apperror"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/dto"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service service.AuthService
}

func NewAuthController(service service.AuthService) *AuthController {
	return &AuthController{service: service}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, apperror.ErrBadRequest)
		return
	}

	session, ok := c.service.Login(req.StudentID, req.Password)
	if !ok {
		err := apperror.ErrUnauthorized
		err.Message = "Invalid student ID or password"
		dto.SendError(ctx, err)
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, dto.NewLoginResponse(session, req.RememberMe))
}
