package middleware

import (
	"WatchHive/pkg/config"
	"WatchHive/pkg/utils/errmsg"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func UserAuthMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{errmsg.MsgErr: errmsg.MsgTokenMissingErr})
		c.Abort()
		return
	}

	cfg, _ := config.LoadConfig()
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.UserAccessKey), nil
	})

	if err != nil || !token.Valid {
		log.Println("Token error:", err)
		c.JSON(http.StatusUnauthorized, gin.H{errmsg.MsgErr: errmsg.MsgTokenErr})
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{errmsg.MsgErr: errmsg.MsgUnAuthErr})
		c.Abort()
		return
	}

	role, ok := claims["role"].(string)
	if !ok || role != "client" {
		c.JSON(http.StatusForbidden, gin.H{errmsg.MsgErr: errmsg.MsgUnAuthErr})
		c.Abort()
		return
	}

	id, ok := claims["id"].(float64)
	if !ok || id == 0 {
		c.JSON(http.StatusForbidden, gin.H{errmsg.MsgErr: errmsg.MsgIdGetErr})
		c.Abort()
		return
	}

	c.Set("role", role)
	c.Set("id", int(id))

	c.Next()
}
