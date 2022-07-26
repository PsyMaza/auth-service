package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"gitlab.com/g6834/team17/auth-service/internal/api/requests"
	"gitlab.com/g6834/team17/auth-service/internal/constants"
	"gitlab.com/g6834/team17/auth-service/internal/interfaces"
	"gitlab.com/g6834/team17/auth-service/internal/models"
	"gitlab.com/g6834/team17/auth-service/internal/utils"
	"net/http"
	"time"
)

type authHandlers struct {
	logger      *zerolog.Logger
	presenters  interfaces.Presenters
	authService interfaces.AuthService
}

func newAuthHandlers(logger *zerolog.Logger, presenter interfaces.Presenters, authService interfaces.AuthService) *authHandlers {
	return &authHandlers{
		logger:      logger,
		presenters:  presenter,
		authService: authService,
	}
}

func AuthRouter(logger *zerolog.Logger, presenter interfaces.Presenters, authService interfaces.AuthService) http.Handler {
	handlers := newAuthHandlers(logger, presenter, authService)

	r := chi.NewRouter()
	r.Get("/all", handlers.all)
	r.Post("/login", handlers.login)
	r.Post("/logout", handlers.logout)
	r.Post("/validate", handlers.validate)

	return r
}

// Login
// @ID login
// @tags auth
// @Summary Authorized user
// @Description Authenticate and authorized user. Return access and refresh tokens in cookies.
// @Accept json
// @Produce json
// @Param redirect_uri query string false "redirect uri"
// @Param login body requests.Login true "request body"
// @Success 200 {object} response.TokenPair true "ok"
// @Header 200 {string} access_token	"token for access services"
// @Header 200 {string} refresh_token	"token for refresh access_token"
// @Failure 400 {object} response.Error "bad request"
// @Failure 404 {string} string "404 page not found"
// @Failure 403 {object} response.Error "forbidden"
// @Failure 500 {object} response.Error "internal error"
// @Router /login [post]
func (handlers *authHandlers) login(w http.ResponseWriter, r *http.Request) {
	ctx, span := utils.StartSpan(r.Context())
	defer span.End()

	var input requests.Login
	err := utils.ReadJson(r, &input)
	if err != nil {
		handlers.presenters.Error(w, r, models.ErrorInternal(err))
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
		handlers.presenters.Error(w, r, models.ErrorForbidden(err))
		return
	}

	atCookie := http.Cookie{
		Name:    constants.ACCESS_TOKEN,
		Value:   td.AccessToken,
		Path:    "/",
		Expires: td.AtExpires,
	}
	rtCookie := http.Cookie{
		Name:     constants.REFRESH_TOKEN,
		Value:    td.RefreshToken,
		Path:     "/",
		Expires:  td.RtExpires,
		HttpOnly: true,
	}

	http.SetCookie(w, &atCookie)
	http.SetCookie(w, &rtCookie)

	redirectUrl := r.URL.Query().Get("redirect_uri")

	if len(redirectUrl) > 0 {
		http.Redirect(w, r, redirectUrl, http.StatusFound)
	} else {
		handlers.presenters.JSON(w, r, td)
	}
}

// Logout
// @ID logout
// @tags auth
// @Summary Clears tokens
// @Description Clears access and refresh tokens
// @Security Auth
// @Produce json
// @Param redirect_uri query string false "redirect uri"
// @Param access_token header string true "access token"
// @Param refresh_token header string true "refresh token"
// @Success 200  "ok"
// @Failure 302  "redirect"
// @Failure 500  "internal error"
// @Router /logout [post]
func (handlers *authHandlers) logout(w http.ResponseWriter, r *http.Request) {
	_, span := utils.StartSpan(r.Context())
	defer span.End()

	rtCookie := http.Cookie{
		Name:     constants.REFRESH_TOKEN,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		Secure:   true,
		HttpOnly: true,
	}

	atCookie := http.Cookie{
		Name:    constants.ACCESS_TOKEN,
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}

	http.SetCookie(w, &atCookie)
	http.SetCookie(w, &rtCookie)

	redirectUrl := r.URL.Query().Get(constants.REDIRECT_URI)

	if len(redirectUrl) > 0 {
		http.Redirect(w, r, redirectUrl, http.StatusFound)
	}

}

// Validate
// @ID Validate
// @tags auth
// @Summary Validate tokens
// @Description Validate tokens and refresh tokens if refresh token is valid
// @Security Auth
// @Produce json
// @Param access_token header string true "access token"
// @Param refresh_token header string true "refresh token"
// @Success 200 {object} response.TokenPair true "ok"
// @Failure 403 {string} string "forbidden"
// @Failure 500 {string} string "internal error"
// @Router /validate [post]
func (handlers *authHandlers) validate(w http.ResponseWriter, r *http.Request) {
	ctx, span := utils.StartSpan(r.Context())
	defer span.End()

	at, err := r.Cookie(constants.ACCESS_TOKEN)
	if err != nil {
		handlers.presenters.Error(w, r, models.ErrorForbidden(err))
		return
	}

	rt, err := r.Cookie(constants.REFRESH_TOKEN)
	if err != nil {
		handlers.presenters.Error(w, r, models.ErrorForbidden(err))
		return
	}

	_, err = handlers.authService.VerifyToken(ctx, &models.TokenPair{
		AccessToken:  at.Value,
		RefreshToken: rt.Value,
	})
	if err != nil {
		handlers.presenters.Error(w, r, models.ErrorForbidden(err))
		return
	}
}

func (handlers *authHandlers) all(w http.ResponseWriter, r *http.Request) {

}
