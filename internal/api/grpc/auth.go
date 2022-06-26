package grpc

import (
	"context"
	"gitlab.com/g6834/team17/api/pkg/auth_service"
	"gitlab.com/g6834/team17/auth-service/internal/interfaces"
	"gitlab.com/g6834/team17/auth-service/internal/models"
)

type AuthApi struct {
	authS interfaces.AuthService
	auth_service.UnimplementedAuthServiceServer
}

func NewAuthAPI(authS interfaces.AuthService) *AuthApi {
	return &AuthApi{authS: authS}
}

func (a *AuthApi) Validate(ctx context.Context, req *auth_service.ValidateTokenRequest) (*auth_service.ValidateTokenResponse, error) {
	tokens, err := a.authS.VerifyToken(ctx, &models.TokenPair{
		AccessToken:  req.AccessToken,
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		return &auth_service.ValidateTokenResponse{
			AccessToken:  "",
			RefreshToken: "",
			Status:       auth_service.Statuses_invalid,
		}, err
	}

	return &auth_service.ValidateTokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		Status:       auth_service.Statuses_valid,
	}, nil
}
