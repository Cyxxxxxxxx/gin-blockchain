package services

import (
	"blockchain/mysql"
)

type Response struct {
	Code int64       `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func JudgeConfirmPermission(userid string) (int, error) {
	//判断数据表中是否有该用户
	count, err := mysql.SelectAuth(userid)
	return count, err
}
