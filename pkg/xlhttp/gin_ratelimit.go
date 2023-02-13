/**
gin 限流中间件
*/
package xlhttp

import (
	"github.com/didip/tollbooth"
	"github.com/gin-gonic/gin"
	"github.com/seanlan/xlvein/pkg/xlerror"
)

// RateLimitMiddleware
// @description gin限流中间件
// @param max float64 每秒限制请求量（QPS）
func RateLimitMiddleware(max float64) gin.HandlerFunc {
	lmt := tollbooth.NewLimiter(max, nil)
	return func(c *gin.Context) {
		r := Build(c)
		httpError := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if httpError != nil {
			r.JsonReturn(xlerror.ErrRateLimit)
			r.ctx.Abort()
		} else {
			c.Next()
		}
	}
}
