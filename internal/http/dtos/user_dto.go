package dtos

import dto_utils "github.com/ViitoJooj/ward/internal/http/dtos/utils"

type AdminCreateUserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Active   *bool  `json:"active"`
}

type AdminUpdateUserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Active   bool   `json:"active"`
}

type SelfUpdateUserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AdminCreateUserData struct {
	User              dto_utils.UserData `json:"user"`
	TemporaryPassword string             `json:"temporary_password"`
}

type AdminCreateUserOutput struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    AdminCreateUserData `json:"data"`
}

type UserOutput struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Data    dto_utils.UserData `json:"data"`
}

type UserListOutput struct {
	Success bool                 `json:"success"`
	Message string               `json:"message"`
	Data    []dto_utils.UserData `json:"data"`
}
