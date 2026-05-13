package route

import (
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/controller"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, c *controller.AuthController) {
	rg.POST("/login", c.Login)
}
