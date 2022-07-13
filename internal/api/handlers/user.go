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
)

type userHandlers struct {
	logger      *zerolog.Logger
	presenters  interfaces.Presenters
	userService interfaces.UserService
}

func newUserHandlers(logger *zerolog.Logger, presenter interfaces.Presenters, userService interfaces.UserService) *userHandlers {
	return &userHandlers{
		logger:      logger,
		presenters:  presenter,
		userService: userService,
	}
}

func UserRouter(logger *zerolog.Logger, presenter interfaces.Presenters, userService interfaces.UserService) http.Handler {
	handlers := newUserHandlers(logger, presenter, userService)

	r := chi.NewRouter()
	r.Get("/i", handlers.get)
	r.Post("/create", handlers.create)

	return r
}

// Create
// @ID create
// @tags user
// @Summary Created new user
// @Description Created new user
// @Accept json
// @Produce json
// @Param access_token header string true "access token"
// @Param refresh_token header string true "refresh token"
// @Param user body requests.CreateUser true "request body"
// @Success 200 {string} string "ok"
// @Failure 400 {object} response.Error "bad request"
// @Failure 403 {object} response.Error "forbidden"
// @Failure 500 {object} response.Error "internal error"
// @Router /create [post]
func (handlers *userHandlers) create(w http.ResponseWriter, r *http.Request) {
	ctx, span := utils.StartSpan(r.Context())
	defer span.End()

	var input requests.CreateUser
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

	err = handlers.userService.Create(ctx, &models.User{
		Username:  input.Username,
		Password:  input.Password,
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	})
	if err != nil {
		handlers.presenters.Error(w, r, err)
		return
	}
}

// GetUserInfo
// @ID GetUserInfo
// @tags user
// @Summary Get user info
// @Description returning user info
// @Accept json
// @Produce json
// @Param access_token header string true "access token"
// @Param refresh_token header string true "refresh token"
// @Success 200 {object} response.User "ok"
// @Failure 400 {object} response.Error "bad request"
// @Failure 403 {object} response.Error "forbidden"
// @Failure 500 {object} response.Error "internal error"
// @Router /i [post]
func (handlers *userHandlers) get(w http.ResponseWriter, r *http.Request) {
	ctx, span := utils.StartSpan(r.Context())
	defer span.End()

	user := ctx.Value(constants.CTX_USER).(*models.User)

	handlers.presenters.JSON(w, r, user)
}
