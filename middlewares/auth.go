package middlewares

import (
	"blockchain/controller"
	"blockchain/services"
	"blockchain/setting"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Resp struct {
	Data string `json:"data"`
	Msg  string `json:"msg"`
}

//判断管理后台是否有超管权限
func JudgePermissionWeb(c *gin.Context) {

	//decode := c.Query("userid")
	//
	//if decode == "" {
	//	controller.ResponseError(c, controller.CodeErrPermission)
	//	c.Abort()
	//	return
	//}
	//
	////解密userid
	//useridInt, err := tool.Decrypt(decode)
	//
	//if err != nil {
	//	zap.L().Error("解密发生错误", zap.Error(err))
	//}
	//userid := fmt.Sprint(useridInt)

	userid := c.Query("userid")

	if userid == "" {
		// controller.ResponseError(c, controller.CodeErrPermission)
		c.Abort()
		return
	}

	//获得授权工号列表
	authArray := strings.Split(setting.Conf.AuthID, ";")

	//循环查找是否有超管权限
	for _, value := range authArray {
		if value == userid {
			c.Next()
			return
		}
	}

	// controller.ResponseError(c, controller.CodeErrPermission)
	c.Abort()
	return
}

//判断管理后台是否具有超管或扫码权限
func JudgePermissionAll(c *gin.Context) {
	userid := c.Query("userid")

	if userid == "" {
		// controller.ResponseError(c, controller.CodeErrPermission)
		c.Abort()
		return
	}
	//decode := c.Query("userid")
	//
	//if decode == "" {
	//	controller.ResponseError(c, controller.CodeErrPermission)
	//	c.Abort()
	//	return
	//}
	//
	////解密userid
	//useridInt, err := tool.Decrypt(decode)
	//if err != nil {
	//	zap.L().Error("解密发生错误", zap.Error(err))
	//}
	//userid := fmt.Sprint(useridInt)

	//获得授权工号列表
	authArray := strings.Split(setting.Conf.AuthID, ";")

	//循环查找是否有超管权限
	for _, value := range authArray {
		if value == userid {
			c.Next()
			return
		}
	}

	//判断是否具有扫码权限
	count, err := services.JudgeConfirmPermission(userid)

	if err != nil {
		zap.L().Error("查找扫码权限失败", zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		c.Abort()
		return
	}

	if count >= 1 {
		//存在扫码权限
		c.Next()
	} else {
		// controller.ResponseError(c, controller.CodeErrPermission)
		c.Abort()
	}
}
