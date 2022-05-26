package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/g6834/team17/auth_service/internal/model"
	"gitlab.com/g6834/team17/auth_service/internal/service/user_service"
	"net/http"
)

type UserCrudHandler interface {
	GetAll(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
}

type userCrudHandler struct {
	us user_service.UserService
}

func NewUserCrudHandler(us user_service.UserService) UserCrudHandler {
	return &userCrudHandler{
		us,
	}
}

func (uc *userCrudHandler) GetAll(ctx *gin.Context) {
	users, err := uc.us.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (uc *userCrudHandler) Create(ctx *gin.Context) {
	var u model.User
	if err := ctx.ShouldBindJSON(&u); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, "Invalid json provider")
		return
	}

	err := uc.us.Create(ctx, &u)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (uc *userCrudHandler) Update(ctx *gin.Context) {
	var u model.User
	if err := ctx.ShouldBindJSON(&u); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, "Invalid json provider")
		return
	}

	err := uc.us.Update(ctx, &u)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}

	ctx.Status(http.StatusOK)
}
