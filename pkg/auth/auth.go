package auth

import (
	"ops/pkg/ctl"
	"net/http"
	"ops/pkg/db"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

func CheckAuth(username, password string) bool {
	var auth db.Auth
	db.Db.Where("`user`=?  and state!=0", username).Get(&auth)

	if db.Aes.Decode(auth.Password) == password {
		ctl.Debug(auth.AuthId)
		return true
	}
	return false
}

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func GetAuth(c *gin.Context) {
	appG := Gin{C: c}
	username := c.PostForm("username")
	password := c.PostForm("password")
	ctl.Debug(username, password)
	valid := validation.Validation{}
	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)
	if !ok {
		MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, INVALID_PARAMS, nil)
		return
	}
	isExist := CheckAuth(username, password)
	if !isExist {
		appG.Response(http.StatusUnauthorized, ERROR_AUTH, nil)
		return
	}
	token, err := GenerateToken(username, password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, ERROR_AUTH_TOKEN, nil)
	}

	c.SetCookie("token", token, 40240, "/", "", false, false)

	c.Redirect(http.StatusFound, "/v1/index")

}
