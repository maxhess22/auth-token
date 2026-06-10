package middleware

import (
	"fmt"
	"max/auth/utils"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Obtener cabecera Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.RespondJSON(c, http.StatusUnauthorized, "Falta el token de autorización", nil)
			c.Abort() // Detiene la petición aquí
			return
		}

		// 2. Extraer el token de "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 3. Parsear y validar el token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validar el método de firma
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			utils.RespondJSON(c, http.StatusUnauthorized, "Token inválido o expirado", nil)
			c.Abort()
			return
		}

		// 4. Extraer los claims y guardarlos en el contexto
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Guardamos el ID del usuario en el contexto de Gin para que el controlador pueda usarlo
			c.Set("userID", claims["sub"])
			c.Next() // Permite que la petición continúe al controlador protegido
		} else {
			utils.RespondJSON(c, http.StatusUnauthorized, "Error leyendo el token", nil)
			c.Abort()
		}
	}
}
