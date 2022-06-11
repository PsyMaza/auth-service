package api

//
//func Login(c *gin.Context) {
//	var u User
//	if err := c.ShouldBindJSON(&u); err != nil {
//		c.JSON(http.StatusUnprocessableEntity, "Invalid json provider")
//		return
//	}
//
//	as := app - service.New("fdsafasdfasd")
//	token, err := as.Authorize(u.Username, u.Password)
//	if err != nil {
//		c.JSON(http.StatusUnauthorized, err)
//		return
//	}
//
//	c.JSON(http.StatusOK, token)
//}
