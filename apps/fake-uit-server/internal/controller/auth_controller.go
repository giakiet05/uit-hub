package controller

import (
	"net/http"

	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/apperror"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/dto"
<<<<<<< HEAD
	"github.com/gin-gonic/gin"
)

func (c *UitController) Login(ctx *gin.Context) {

=======
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
>>>>>>> origin/dev
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, apperror.ErrBadRequest)
		return
	}

<<<<<<< HEAD
	res, appErr := c.service.Login(req)
	if appErr != nil {
		dto.SendError(ctx, *appErr)
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, res.Data)
=======
	session, ok := c.service.Login(req.StudentID, req.Password)
	if !ok {
		err := apperror.ErrUnauthorized
		err.Message = "Invalid student ID or password"
		dto.SendError(ctx, err)
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, dto.NewLoginResponse(session, req.RememberMe))
>>>>>>> origin/dev
}
