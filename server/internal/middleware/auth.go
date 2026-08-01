package middleware

import (
	"net/http"
	"strings"

	"git-backup-web/server/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// JWTAuth 校验 Authorization: Bearer <token>
func JWTAuth(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			return
		}
		tokenStr := header[7:]
		var cfg config.Config
		database.First(&cfg)
		secret := cfg.JWTSecret
		if secret == "" {
			secret = "git-backup-default-secret"
		}
		parsed, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !parsed.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "登录已失效"})
			return
		}
		if claims, ok := parsed.Claims.(jwt.MapClaims); ok {
			if name, ok := claims["user"].(string); ok {
				c.Set("username", name)
			}
		}
		c.Next()
	}
}
