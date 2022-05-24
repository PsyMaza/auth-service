package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/g6834/team17/auth_service/internal/service/auth_service"
	"net/http"
)

func Login(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provider")
		return
	}

	as := auth_service.New("fdsafasdfasd")
	token, err := as.Authorize(u.Username, u.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, token)
}
