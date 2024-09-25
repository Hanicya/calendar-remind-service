package validate

import (
	"calendar-remind-service/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 密钥
var jwtKey []byte = []byte{87, 50, 212, 174, 174, 20, 169, 128, 174, 65, 243, 200, 149, 127, 66, 190, 184, 59, 49, 124, 151, 214, 51, 166, 184, 193, 9, 175, 69, 40, 205, 38}

// TokenMiddleware 令牌校验中间件
func TokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查请求路径，如果是 WebSocket 路由，则跳过验证
		if c.Request.URL.Path == "/ws" {
			c.Next()
			return
		}
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}
		issuer, err := util.ParseJwt(jwtKey, token)
		// token 的验证
		if issuer == -1 || err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		// 存储在上下文中
		c.Set("userID", issuer)
		// 继续处理请求
		c.Next()
	}
}
