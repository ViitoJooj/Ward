package dto_utils

type UserData struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Active     bool   `json:"active"`
	Updated_at string `json:"updated_at"`
	Created_at string `json:"created_at"`
}
