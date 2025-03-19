package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// Secret key untuk JWT
var jwtKey []byte

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to load .env file")
	}
	jwtKey = []byte(os.Getenv("JWT_SECRET"))
}

// GenerateToken - Membuat JWT Token untuk user
func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token berlaku 24 jam
	claims := &jwt.MapClaims{
		"username": username,
		"exp":      expirationTime.Unix(),
	}

	// Buat token dengan algoritma HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// Middleware untuk memvalidasi token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
			c.Abort()
			return
		}

		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			// gk bisa di bagian token dgn error "error": "token is malformed: could not base64 decode header: illegal base64 data at input byte 6"
			// update : isi header authorization tidak boleh ada bearer
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token invalid"})
			c.Abort()
			return
		}

		// Simpan username ke context agar bisa diakses di handler selanjutnya
		c.Set("username", (*claims)["username"])
		c.Next()
	}
}
