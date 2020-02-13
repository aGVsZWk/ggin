package api

import (
	"github.com/gin-gonic/gin"
	"github.com/astaxie/beego/validation"
	"ggin/pkg/e"
	"ggin/models"
	"ggin/pkg/util"
	"net/http"
	"ggin/pkg/app"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func GetAuth(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}

	username := c.Query("username")
	password := c.Query("password")

	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	if !ok {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	isExist, err := models.CheckAuth(username, password)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	if !isExist {
		appG.Response(http.StatusOK, e.ERROR_AUTH, nil)
		return
	}
	token, err := util.GenerateToken(username, password)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_AUTH_TOKEN, nil)
	}
	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})
}
