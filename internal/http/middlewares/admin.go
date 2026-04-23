package middlewares

import (
	"encoding/json"

	dto_utils "github.com/ViitoJooj/ward/internal/http/dtos/utils"
	"github.com/valyala/fasthttp"
)

func AdminOnlyMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		roleRaw := ctx.UserValue("userRole")
		role, ok := roleRaw.(string)
		if !ok || role != "admin" {
			output := dto_utils.Error{
				Success: false,
				Message: "forbidden",
			}
			res, _ := json.Marshal(output)
			ctx.SetStatusCode(fasthttp.StatusForbidden)
			ctx.SetContentType("application/json")
			ctx.SetBody(res)
			return
		}

		next(ctx)
	}
}
