package auth

import (
	"ops/pkg/ctl"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, Response{
		Code: httpCode,
		Msg:  GetMsg(errCode),
		Data: data,
	})
	return
}

func BindAndValid(c *gin.Context, form interface{}) (int, int) {
	err := c.Bind(form)
	if err != nil {
		return http.StatusBadRequest, INVALID_PARAMS
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return http.StatusInternalServerError, ERROR
	}
	if !check {
		MarkErrors(valid.Errors)
		return http.StatusBadRequest, INVALID_PARAMS
	}

	return http.StatusOK, SUCCESS
}

func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		ctl.Log.Info(err.Key, err.Message)
	}
	return
}
