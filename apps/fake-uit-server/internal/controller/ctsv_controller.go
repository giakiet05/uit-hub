package controller

import (
	"net/http"

	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/apperror"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/dto"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/middleware"
	"github.com/gin-gonic/gin"
)

func (c *UitController) ConfirmLetter(ctx *gin.Context) {

	studentID, ok := ctx.Get(middleware.StudentIDKey)
	if !ok {
		dto.SendError(ctx, apperror.ErrUnauthorized)
		return
	}

	var req dto.ConfirmLetterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, apperror.ErrBadRequest)
		return
	}

	res, appErr := c.service.ConfirmLetter(studentID.(string), req)
	if appErr != nil {
		dto.SendError(ctx, *appErr)
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, res.Data)
}

func (c *UitController) BankLoans(ctx *gin.Context) {

	studentID, ok := ctx.Get(middleware.StudentIDKey)
	if !ok {
		dto.SendError(ctx, apperror.ErrUnauthorized)
		return
	}

	var req dto.BankLoansRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, apperror.ErrBadRequest)
		return
	}

	res, appErr := c.service.BankLoans(studentID.(string), req)
	if appErr != nil {
		dto.SendError(ctx, *appErr)
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, res.Data)
}

func (c *UitController) TrainingPointConfirm(ctx *gin.Context) {

	studentID, ok := ctx.Get(middleware.StudentIDKey)
	if !ok {
		dto.SendError(ctx, apperror.ErrUnauthorized)
		return
	}

	var req dto.TrainingPointConfirmRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, apperror.ErrBadRequest)
		return
	}

	res, appErr := c.service.TrainingPointConfirm(studentID.(string), req)
	if appErr != nil {
		dto.SendError(ctx, *appErr)
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, res.Data)
}

func (c *UitController) LanguageCertificate(ctx *gin.Context) {

	studentID, ok := ctx.Get(middleware.StudentIDKey)
	if !ok {
		dto.SendError(ctx, apperror.ErrUnauthorized)
		return
	}

	var req dto.LanguageCertificateUploadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, apperror.ErrBadRequest)
		return
	}

	res, appErr := c.service.LanguageCertificate(studentID.(string), req)
	if appErr != nil {
		dto.SendError(ctx, *appErr)
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, res.Data)
}
