package middleware

import (
	"ecommerce-api/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware adalah satpam pemeriksa tiket JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil Header "Authorization" dari request
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak: silakan sertakan token JWT!"})
			return
		}

		// 2. Standar industri: Token harus berformat "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Format token salah, gunakan: Bearer <token>"})
			return
		}

		tokenString := parts[1]

		// 3. Periksa keaslian stempel token menggunakan validator kita
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token palsu atau sudah kadaluarsa!"})
			return
		}

		// 4. Titipkan identitas user ke context (bisa dipakai handler berikutnya)
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		// 5. TIKET ASLI! Silakan lewat ke ruangan handler
		c.Next()
	}
}
