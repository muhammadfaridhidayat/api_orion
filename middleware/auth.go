package middleware

import (
	"api_orion/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		// Try Authorization header first (Bearer token), then fall back to cookie
		var tokenString string

		authHeader := ctx.GetHeader("Authorization")
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			cookie, err := ctx.Cookie("session_token")
			if err != nil {
				ctx.JSON(http.StatusUnauthorized, model.ErrorResponse{
					Success: false,
					Status:  http.StatusUnauthorized,
					Message: "Error from auth middleware",
					Errors: map[string]string{
						"error": "unauthorized",
					},
				})
				ctx.Abort()
				return
			}
			tokenString = cookie
		}

		claims := &model.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return model.JwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				ctx.JSON(http.StatusUnauthorized, model.ErrorResponse{
					Success: false,
					Status:  http.StatusUnauthorized,
					Message: "Error from auth middleware",
					Errors: map[string]string{
						"error": "unauthorized",
					},
				})
			} else {
				ctx.JSON(http.StatusBadRequest, model.ErrorResponse{
					Success: false,
					Status:  http.StatusBadRequest,
					Message: "Error from auth middleware",
					Errors: map[string]string{
						"error": "bad request",
					},
				})
			}
			ctx.Abort()
			return
		}

		if !token.Valid {
			ctx.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Success: false,
				Status:  http.StatusUnauthorized,
				Message: "Error from auth middleware",
				Errors: map[string]string{
					"error": "unauthorized",
				},
			})
			ctx.Abort()
			return
		}

		ctx.Set("id", claims.UserID)

		ctx.Next()
	})
}
