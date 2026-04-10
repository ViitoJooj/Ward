package middlewares

import (
	"log"

	"github.com/ViitoJooj/door/pkg/jwtTokens"
	"github.com/golang-jwt/jwt/v4"
	"github.com/valyala/fasthttp"
)

func UserIdMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		tokenString := string(ctx.Request.Header.Cookie("token"))
		if tokenString == "" {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			return
		}

		token, err := jwtTokens.ValidateToken(tokenString)
		if err != nil {
			log.Println(err)
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("invalid token claims")
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			return
		}

		userIdFloat, ok := claims["user_id"].(float64)
		if !ok {
			log.Println("user_id not found or invalid type")
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			return
		}

		userId := int64(userIdFloat)

		ctx.SetUserValue("userId", userId)
		next(ctx)
	}
}
