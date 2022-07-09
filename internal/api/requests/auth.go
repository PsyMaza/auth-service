package requests

// swagger:model Login
type Login struct {
	// Username for authentication
	Username string `json:"login" validate:"required" example:"test123"`

	// Password for authentication
	Password string `json:"password" validate:"required" example:"qwerty"`
}
