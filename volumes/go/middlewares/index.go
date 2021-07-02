package middlewares

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	limit "github.com/yangxikun/gin-limit-by-key"
	"golang.org/x/time/rate"
)

// 找出可用的網域
func AllowDomains(domains []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		isAllow := false
		for _, val := range domains {
			if strings.Index(c.Request.Host, val) != -1 {
				isAllow = true
			}
		}
		if isAllow {
			c.Next()
		} else {
			c.AbortWithStatus(401)
		}
	}
}

// 設定每人瞬間的limit request
func SetRequestLimit() gin.HandlerFunc {
	return limit.NewRateLimiter(func(c *gin.Context) string {
		return c.ClientIP() // limit rate by client ip
	}, func(c *gin.Context) (*rate.Limiter, time.Duration) {
		return rate.NewLimiter(rate.Every(100*time.Millisecond), 5), time.Hour // limit 5 qps/clientIp and permit bursts of at most 5 tokens, and the limiter livene
	}, func(c *gin.Context) {
		c.AbortWithStatus(429) // handle exceed rate limit request
	})
}
