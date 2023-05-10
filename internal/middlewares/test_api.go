package middlewares

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"github.com/gerladeno/chat-service/internal/types"
)

func SetToken(c echo.Context, uid types.UserID) {
	claims := claims{}
	claims.StandardClaims.Subject = uid.String()
	c.Set(tokenCtxKey, &jwt.Token{
		Claims: claims,
		Valid:  true,
	})
}
