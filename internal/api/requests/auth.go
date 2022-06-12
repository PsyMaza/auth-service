package requests

type Login struct {
	Username string `json:"username" validate:"required" example:"user123"`
	Password string `json:"password" validate:"required" example:"qwerty"`
}
