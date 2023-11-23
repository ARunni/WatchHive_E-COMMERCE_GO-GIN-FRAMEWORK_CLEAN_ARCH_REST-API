package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt/v4"
)

func UserAUthMiddleware(c *gin.Context) {
	tokenstring := c.GetHeader("Authorization")
	if tokenstring == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization token"})
		c.Abort()
		return
	}
	tokenstring = strings.TrimPrefix(tokenstring, "Bearer")

	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {

		return []byte("watchhive@123"), nil

	})
	fmt.Println("kkkk", token)
	if err != nil || !token.Valid {
		fmt.Println("tttttttttttttt", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization Token"})
		c.Abort()
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		c.Abort()
		return
	}

	role, ok := claims["role"].(string)
	if !ok || role != "client" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access "})
		c.Abort()
		return
	}
	id, ok := claims["id"].(float64)
	if !ok || id == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "error in retrieving id"})
		c.Abort()
		return
	}
	c.Set("role", role)
	c.Set("id", int(id))

	c.Next()
}
