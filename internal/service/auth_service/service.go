package auth

type AuthService struct {
	SecretKey string
}

func New(secretKey string) *AuthService {
	return &AuthService{
		secretKey,
	}
}

var (
	authorized = "authorized"
	userId     = "user_id"
	exp        = "exp"
)
