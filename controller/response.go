package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
{
	"code": //程序中的错误码
	"msg": //提示信息
	“data" : //数据
}

*/

type ResponeseCode struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}
type ResponeseCode02 struct {
	Code   ResCode     `json:"code"`
	Msg    interface{} `json:"msg"`
	Result interface{} `json:"result"`
	Data   interface{} `json:"data"`
}

// ResponseError 返回一个code 类型的错误
func ResponseError(c *gin.Context, code ResCode) {
	c.JSON(http.StatusOK, &ResponeseCode{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

// ResponseErrorWithMsg 返回一个提示信息为自定义的错误
func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	c.JSON(http.StatusOK, &ResponeseCode{
		Code: code,
		Msg:  code.Msg(),
		Data: msg,
	})
}

// ResponseSuccess 返回成功的信息
func ResponseSuccess(c *gin.Context, code ResCode, data interface{}) {
	c.JSON(http.StatusOK, &ResponeseCode{
		Code: code,
		Msg:  code.Msg(),
		Data: data,
	})
}

// ResponseSuccess01 返回系统的状态的信息
func ResponseMsg(c *gin.Context, code ResCode) {
	c.JSON(http.StatusOK, &ResponeseCode{
		Code: code,
		Msg:  code.Msg(),
	})
}

// ResponseSuccessInput 返回导入成功的信息
func ResponseSuccessInput(c *gin.Context, code ResCode, data02 interface{}, data01 interface{}) {
	c.JSON(http.StatusOK, &ResponeseCode02{
		Code:   code,
		Msg:    code.Msg(),
		Result: data02,
		Data:   data01,
	})
}
