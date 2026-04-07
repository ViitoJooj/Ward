package dtos

// Utils
type UserData struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Updated_at string `json:"updated_at"`
	Created_at string `json:"created_at"`
}

// Register
type RegisterInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterOutput struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    UserData `json:"data"`
}

// Login
type LoginInput struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type LoginOutput struct {
	Success      bool     `json:"success"`
	Message      string   `json:"message"`
	Data         UserData `json:"data"`
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
}
