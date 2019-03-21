package middlewares

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
)

const (
	ContextGetRedisName = "GETREDIS"
	ContextSetRedisName = "SETREDIS"
	/* Redis Keys */
	// 무료 매칭 가능  2회 Zincrby
	FREEMATCH_COUNT = "freematch"

	MATCH_RECOMMAND = "matchRecommand"

	TOTALMATCHCOUNT = "totalmatchcount"
)

func ContextRedis(service string, client *redis.Client) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			c.SetRequest(req.WithContext(context.WithValue(req.Context(), service, client)))
			return next(c)
		}
	}
}
