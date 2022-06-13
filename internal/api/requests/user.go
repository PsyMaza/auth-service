package requests

type CreateUser struct {
	Username  string `json:"username" validate:"required" example:"user123"`
	Password  string `json:"password" validate:"required" example:"qwerty"`
	Email     string `json:"email" validate:"required" example:"qwerty@domen.com"`
	FirstName string `json:"first_name" validate:"required" example:"qwerty"`
	LastName  string `json:"last_name" validate:"required" example:"qwerty"`
}
