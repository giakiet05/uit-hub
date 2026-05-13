package middleware

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/apperror"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/config"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/dto"
	"github.com/gin-gonic/gin"
)

func Chaos(cfg config.ChaosConfig) gin.HandlerFunc {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	return func(ctx *gin.Context) {
		if delayMS, ok := forcedDelay(ctx); ok {
			sleep(delayMS)
		} else if cfg.Enabled {
			sleep(randomDelay(cfg, random))
		}

		if status, ok := forcedStatus(ctx); ok {
			abortWithStatus(ctx, status)
			return
		}

		if cfg.Enabled && cfg.ErrorRate > 0 && random.Float64() < normalizedRate(cfg.ErrorRate) {
			abortWithStatus(ctx, http.StatusInternalServerError)
			return
		}

		ctx.Next()
	}
}

func forcedDelay(ctx *gin.Context) (int, bool) {
	delay := ctx.GetHeader("X-Fake-Delay")
	if delay == "" {
		delay = ctx.Query("__fake_delay")
	}
	if delay == "" {
		return 0, false
	}

	delayMS, err := strconv.Atoi(delay)
	if err != nil || delayMS <= 0 {
		return 0, false
	}

	return delayMS, true
}

func forcedStatus(ctx *gin.Context) (int, bool) {
	statusValue := ctx.GetHeader("X-Fake-Status")
	if statusValue == "" {
		statusValue = ctx.Query("__fake_status")
	}
	if statusValue == "" {
		return 0, false
	}

	status, err := strconv.Atoi(statusValue)
	if err != nil || status < 400 || status > 599 {
		return 0, false
	}

	return status, true
}

func randomDelay(cfg config.ChaosConfig, random *rand.Rand) int {
	minDelay := max(cfg.MinDelayMS, 0)
	maxDelay := max(cfg.MaxDelayMS, 0)

	if maxDelay <= minDelay {
		return minDelay
	}

	return random.Intn(maxDelay-minDelay+1) + minDelay
}

func sleep(delayMS int) {
	if delayMS <= 0 {
		return
	}

	time.Sleep(time.Duration(delayMS) * time.Millisecond)
}

func normalizedRate(rate float64) float64 {
	if rate < 0 {
		return 0
	}
	if rate > 1 {
		return 1
	}

	return rate
}

func abortWithStatus(ctx *gin.Context, status int) {
	dto.AbortWithError(ctx, apperror.FromStatus(status))
}
