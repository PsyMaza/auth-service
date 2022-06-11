package auth

type AuthService struct {
	SecretKey string
}

func New(secretKey string) *AuthService {
	t := AuthService{
		SecretKey: 5,
	}
	return &AuthService{
		secretKey,
	}
}

var (
	authorized = "authorized"
	userId     = "user_id"
	exp        = "exp"
)
