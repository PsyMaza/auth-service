package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"gitlab.com/g6834/team17/auth-service/internal/api/requests"
	"gitlab.com/g6834/team17/auth-service/internal/interfaces"
	"gitlab.com/g6834/team17/auth-service/internal/models"
	"gitlab.com/g6834/team17/auth-service/internal/utils"
	"net/http"
)

const (
	ACCESS_TOKEN  = "access_token"
	REFRESH_TOKEN = "refresh_token"
)

type authHandlers struct {
	logger      zerolog.Logger
	presenters  interfaces.Presenters
	authService interfaces.AuthService
}

func newAuthHandlers(logger zerolog.Logger, presenter interfaces.Presenters, authService interfaces.AuthService) *authHandlers {
	return &authHandlers{
		logger:      logger,
		presenters:  presenter,
		authService: authService,
	}
}

func AuthRouter(logger zerolog.Logger, presenter interfaces.Presenters, authService interfaces.AuthService) http.Handler {
	handlers := newAuthHandlers(logger, presenter, authService)

	r := chi.NewRouter()
	r.Post("/login", handlers.login)

	return r
}

func (handlers *authHandlers) login(w http.ResponseWriter, r *http.Request) {
	ctx, span := utils.StartSpan(r.Context())
	defer span.End()

	var input requests.Login
	err := utils.ReadJson(r, &input)
	if err != nil {
		handlers.presenters.Error(w, r, models.ErrorBadRequest(err))
		return
	}

	v := validator.New()
	err = v.Struct(input)
	if err != nil {
		handlers.presenters.Error(w, r, models.ErrorBadRequest(err))
		return
	}

	td, err := handlers.authService.Authorize(ctx, input.Username, input.Password)
	if err != nil {
		handlers.presenters.Error(w, r, err)
		return
	}

	atCookie := http.Cookie{
		Name:    ACCESS_TOKEN,
		Value:   td.AccessToken,
		Path:    "/",
		Expires: td.AtExpires,
		Secure:  false,
	}
	rtCookie := http.Cookie{
		Name:     REFRESH_TOKEN,
		Value:    td.RefreshToken,
		Path:     "/",
		Expires:  td.RtExpires,
		Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(w, &rtCookie)
	http.SetCookie(w, &atCookie)

	handlers.presenters.JSON(w, r, td)
}