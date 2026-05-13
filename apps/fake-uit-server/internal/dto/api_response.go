package dto

import (
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/apperror"
	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ErrorResponse struct {
	Message   string `json:"message"`
	ErrorCode string `json:"error_code"`
}

func SendSuccess(ctx *gin.Context, status int, data any) {
	ctx.JSON(status, SuccessResponse{
		Message: "Successfully!",
		Data:    data,
	})
}

func SendError(ctx *gin.Context, err apperror.AppError) {
	ctx.JSON(err.Status, ErrorResponse{
		Message:   err.Message,
		ErrorCode: err.Code,
	})
}

func AbortWithError(ctx *gin.Context, err apperror.AppError) {
	ctx.AbortWithStatusJSON(err.Status, ErrorResponse{
		Message:   err.Message,
		ErrorCode: err.Code,
	})
}
