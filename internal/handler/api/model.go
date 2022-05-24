package api

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password.go"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
