package handler

import (
	"encoding/json"

	"github.com/ViitoJooj/door/internal/http/dtos"
	"github.com/ViitoJooj/door/internal/services"
	"github.com/valyala/fasthttp"
)

type ApplicationHandler struct {
	applicationService *services.ApplicationService
}

func NewApplicationHandler(applicationService *services.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{
		applicationService: applicationService,
	}
}

func (h *ApplicationHandler) Create(ctx *fasthttp.RequestCtx) {
	var input dtos.ApplicationInput

	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody([]byte(`{"error":"invalid json"}`))
		return
	}

	userIdRaw := ctx.UserValue("userId")
	if userIdRaw == nil {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return
	}

	userId, ok := userIdRaw.(int64)
	if !ok {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return
	}

	createApplication, user, err := h.applicationService.Create(input.Url, input.Country, userId)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadGateway)
		return
	}

	output := dtos.ApplicationOutput{
		Success: true,
		Message: "Application has created with successfull.",
		Created_by: dtos.UserData{
			ID:         user.ID,
			Username:   user.Username,
			Email:      user.Email,
			Updated_at: user.Updated_at.String(),
			Created_at: user.Created_at.String(),
		},
		Data: dtos.ApplicationData{
			ID:      createApplication.ID,
			Url:     createApplication.Url,
			Country: createApplication.Country,
		},
	}

	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}
