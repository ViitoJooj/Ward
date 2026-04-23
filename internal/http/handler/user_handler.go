package handler

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/ViitoJooj/ward/internal/domain"
	"github.com/ViitoJooj/ward/internal/http/dtos"
	dto_utils "github.com/ViitoJooj/ward/internal/http/dtos/utils"
	"github.com/ViitoJooj/ward/internal/services"
	"github.com/valyala/fasthttp"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Create(ctx *fasthttp.RequestCtx) {
	var input dtos.AdminCreateUserInput
	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		output := dto_utils.Error{
			Success: false,
			Message: "invalid json.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	userID, ok := getCurrentUserID(ctx)
	if !ok {
		output := dto_utils.Error{
			Success: false,
			Message: "internal error.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	active := true
	if input.Active != nil {
		active = *input.Active
	}

	user, temporaryPassword, err := h.userService.CreateByAdmin(userID, input.Username, input.Email, input.Role, active)
	if err != nil {
		if errors.Is(err, services.ErrUserForbidden) {
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

		output := dto_utils.Error{
			Success: false,
			Message: err.Error(),
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	output := dtos.AdminCreateUserOutput{
		Success: true,
		Message: "user created.",
		Data: dtos.AdminCreateUserData{
			User:              mapUserData(user),
			TemporaryPassword: temporaryPassword,
		},
	}

	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *UserHandler) GetAll(ctx *fasthttp.RequestCtx) {
	userID, ok := getCurrentUserID(ctx)
	if !ok {
		output := dto_utils.Error{
			Success: false,
			Message: "internal error.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	users, err := h.userService.GetAll(userID)
	if err != nil {
		if errors.Is(err, services.ErrUserForbidden) {
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

		output := dto_utils.Error{
			Success: false,
			Message: "internal error.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadGateway)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	data := make([]dto_utils.UserData, 0, len(users))
	for _, user := range users {
		data = append(data, mapUserData(user))
	}

	output := dtos.UserListOutput{
		Success: true,
		Message: "users fetched.",
		Data:    data,
	}

	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *UserHandler) GetByID(ctx *fasthttp.RequestCtx) {
	userID, ok := getCurrentUserID(ctx)
	if !ok {
		output := dto_utils.Error{
			Success: false,
			Message: "internal error.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	targetID, ok := getPathUserID(ctx)
	if !ok {
		output := dto_utils.Error{
			Success: false,
			Message: "invalid id.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	user, err := h.userService.GetByID(userID, targetID)
	if err != nil {
		if errors.Is(err, services.ErrUserForbidden) {
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

		if errors.Is(err, services.ErrUserNotFound) {
			output := dto_utils.Error{
				Success: false,
				Message: "user not found",
			}
			res, _ := json.Marshal(output)
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.SetContentType("application/json")
			ctx.SetBody(res)
			return
		}

		output := dto_utils.Error{
			Success: false,
			Message: "internal error.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadGateway)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	output := dtos.UserOutput{
		Success: true,
		Message: "user fetched.",
		Data:    mapUserData(user),
	}

	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *UserHandler) UpdateByID(ctx *fasthttp.RequestCtx) {
	var input dtos.AdminUpdateUserInput
	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		output := dto_utils.Error{
			Success: false,
			Message: "invalid json.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	userID, ok := getCurrentUserID(ctx)
	if !ok {
		output := dto_utils.Error{
			Success: false,
			Message: "internal error.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	targetID, ok := getPathUserID(ctx)
	if !ok {
		output := dto_utils.Error{
			Success: false,
			Message: "invalid id.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	user, err := h.userService.UpdateByAdmin(userID, targetID, input.Username, input.Email, input.Password, input.Role, input.Active)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			output := dto_utils.Error{
				Success: false,
				Message: "user not found",
			}
			res, _ := json.Marshal(output)
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.SetContentType("application/json")
			ctx.SetBody(res)
			return
		}

		if errors.Is(err, services.ErrUserForbidden) {
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

		output := dto_utils.Error{
			Success: false,
			Message: err.Error(),
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	output := dtos.UserOutput{
		Success: true,
		Message: "user updated.",
		Data:    mapUserData(user),
	}

	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *UserHandler) DeleteByID(ctx *fasthttp.RequestCtx) {
	userID, ok := getCurrentUserID(ctx)
	if !ok {
		output := dto_utils.Error{
			Success: false,
			Message: "internal error.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	targetID, ok := getPathUserID(ctx)
	if !ok {
		output := dto_utils.Error{
			Success: false,
			Message: "invalid id.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	err := h.userService.DeleteByID(userID, targetID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			output := dto_utils.Error{
				Success: false,
				Message: "user not found",
			}
			res, _ := json.Marshal(output)
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.SetContentType("application/json")
			ctx.SetBody(res)
			return
		}

		if errors.Is(err, services.ErrUserForbidden) {
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

		output := dto_utils.Error{
			Success: false,
			Message: "internal error.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadGateway)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	output := map[string]any{
		"success": true,
		"message": "user deleted.",
	}

	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func (h *UserHandler) UpdateMe(ctx *fasthttp.RequestCtx) {
	var input dtos.SelfUpdateUserInput
	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		output := dto_utils.Error{
			Success: false,
			Message: "invalid json.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	userID, ok := getCurrentUserID(ctx)
	if !ok {
		output := dto_utils.Error{
			Success: false,
			Message: "internal error.",
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	user, err := h.userService.UpdateOwnData(userID, input.Username, input.Email, input.Password)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			output := dto_utils.Error{
				Success: false,
				Message: "user not found",
			}
			res, _ := json.Marshal(output)
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.SetContentType("application/json")
			ctx.SetBody(res)
			return
		}

		output := dto_utils.Error{
			Success: false,
			Message: err.Error(),
		}
		res, _ := json.Marshal(output)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
		return
	}

	output := dtos.UserOutput{
		Success: true,
		Message: "user updated.",
		Data:    mapUserData(user),
	}

	res, _ := json.Marshal(output)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(res)
}

func getCurrentUserID(ctx *fasthttp.RequestCtx) (int, bool) {
	userIDRaw := ctx.UserValue("userId")
	if userIDRaw == nil {
		return 0, false
	}

	userID, ok := userIDRaw.(int)
	if !ok {
		return 0, false
	}

	return userID, true
}

func getPathUserID(ctx *fasthttp.RequestCtx) (int, bool) {
	pathValue := string(ctx.Path())
	userIDString := strings.TrimPrefix(pathValue, "/ward/api/v1/users/")
	if userIDString == "" {
		return 0, false
	}

	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		return 0, false
	}

	return userID, true
}

func mapUserData(user *domain.User) dto_utils.UserData {
	return dto_utils.UserData{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		Role:       user.Role,
		Active:     user.Active,
		Updated_at: user.Updated_at.String(),
		Created_at: user.Created_at.String(),
	}
}
