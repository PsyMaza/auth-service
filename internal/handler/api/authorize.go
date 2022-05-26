package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/g6834/team17/auth_service/internal/model"
	"gitlab.com/g6834/team17/auth_service/internal/service/auth_service"
	"net/http"
)

type AuthorizeHandler interface {
	Login(ctx *gin.Context)
}

type authHandler struct {
	as auth_service.AuthService
}

func NewAuthorizeHandler(as auth_service.AuthService) AuthorizeHandler {
	return &authHandler{as}
}

func (a *authHandler) Login(ctx *gin.Context) {
	var u model.User
	if err := ctx.ShouldBindJSON(&u); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, "Invalid json provider")
		return
	}

	token, err := a.as.Authorize(ctx, u.Username, u.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}

	ctx.JSON(http.StatusOK, token)
}
