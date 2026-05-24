package controller

import (
	"net/http"

	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/apperror"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/dto"
	"github.com/gin-gonic/gin"
)

func (c *UitController) GetRoomsAvailability(ctx *gin.Context) {

	var query dto.QueryRoomsAvailability
	if err := ctx.ShouldBindQuery(&query); err != nil {
		dto.SendError(ctx, apperror.ErrBadRequest)
		return
	}

	res, appErr := c.service.GetRoomsAvailability(query)
	if appErr != nil {
		dto.SendError(ctx, *appErr)
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, res.Data)
}
