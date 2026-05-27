package controller

import (
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/service"
)

type UitController struct {
	service *service.Service
}

func NewUitController(service *service.Service) *UitController {
	return &UitController{service: service}
}
