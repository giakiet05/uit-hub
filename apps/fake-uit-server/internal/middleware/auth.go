package middleware

import (
	"strings"

	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/apperror"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/dto"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/service"
	"github.com/gin-gonic/gin"
)

const StudentIDKey = "student_id"

func RequireAuth(authService service.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, ok := bearerToken(ctx.GetHeader("Authorization"))
		if !ok {
			dto.AbortWithError(ctx, apperror.ErrUnauthorized)
			return
		}

		authUser, ok := authService.ValidateToken(token)
		if !ok {
			dto.AbortWithError(ctx, apperror.ErrUnauthorized)
			return
		}

		ctx.Set(StudentIDKey, authUser.StudentID)
		ctx.Next()
	}
}

func bearerToken(header string) (string, bool) {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return "", false
	}

	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	if token == "" {
		return "", false
	}

	return token, true
}
