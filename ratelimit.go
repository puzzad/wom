package wom

import (
	"github.com/labstack/echo/v5"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"log"
	"net"
	"net/http"
	"time"
)

var rateLimiter = limiter.New(
	memory.NewStore(),
	limiter.Rate{
		Period: time.Minute,
		Limit:  5,
	},
	limiter.WithTrustForwardHeader(true),
	limiter.WithIPv4Mask(net.CIDRMask(28, 32)),
	limiter.WithIPv6Mask(net.CIDRMask(64, 128)),
)

var internalError = echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
var rateLimitError = echo.NewHTTPError(http.StatusTooManyRequests, "Rate limit exceeded")

func rateLimit(r *http.Request) error {
	key := rateLimiter.GetIPKey(r)

	context, err := rateLimiter.Get(r.Context(), key)
	if err != nil {
		log.Printf("Failed to get rate limit context: %v", err)
		return internalError
	}

	if context.Reached {
		log.Printf("Rate limit reached by %s (mask %s)", rateLimiter.GetIP(r), key)
		return rateLimitError
	}

	return nil
}
