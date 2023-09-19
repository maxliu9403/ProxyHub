/*
@Date: 2021/1/12 下午2:24
@Author: max.liu
@File : base
@Desc:
*/

package common

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseController struct{}

type Response struct {
	RetCode RetCode     `json:"RetCode"`
	Message string      `json:"Message"`
	DataSet interface{} `json:"Data"`
}

type ResponseWithTotalCount struct {
	RetCode    RetCode     `json:"RetCode"`
	Message    string      `json:"Message"`
	DataSet    interface{} `json:"Data"`
	TotalCount int64       `json:"TotalCount"`
}

// CheckParams check params, params must be a pointer
func (c *BaseController) CheckParams(ctx *gin.Context, params interface{}) bool {
	code, err := BindAndValid(ctx, params)
	if err != nil {
		c.Response(ctx, nil, NewErrorCode(code, err))
		return false
	}

	return true
}

func (c *BaseController) Response(ctx *gin.Context, data interface{}, err error) {
	httpCode := http.StatusOK

	var msg string
	var retCode RetCode
	if err != nil {
		codeErr, ok := err.(CodeWithErr)
		if ok {
			retCode = codeErr.RetCode
			if codeErr.Error() == "" {
				retCode = SUCCESS
			}
		} else {
			retCode = FAILED
		}
		msg = GetMsg(retCode)
		if len(msg) == 0 {
			msg = codeErr.Error()
		}
		msg = fmt.Sprintf("%s, %s", msg, err.Error())
	} else {
		retCode = SUCCESS
		msg = GetMsg(retCode)
	}

	ctx.JSON(httpCode, Response{
		RetCode: retCode,
		Message: msg,
		DataSet: data,
	})
}

func (c *BaseController) ResponseWithTotalCount(ctx *gin.Context, data interface{}, totalCount int64, err error) {
	httpCode := http.StatusOK
	var msg string
	var retCode RetCode
	if err != nil {
		codeErr, ok := err.(CodeWithErr)
		if ok {
			retCode = codeErr.RetCode
			if codeErr.Error() == "" {
				retCode = SUCCESS
			}
		} else {
			retCode = FAILED
		}
		msg = GetMsg(retCode)
		if len(msg) == 0 {
			msg = codeErr.Error()
		}
		msg = fmt.Sprintf("%s, %s", msg, err.Error())
		ctx.JSON(httpCode, ResponseWithTotalCount{
			RetCode:    retCode,
			Message:    msg,
			DataSet:    data,
			TotalCount: totalCount,
		})
		return
	}

	retCode = SUCCESS
	msg = GetMsg(retCode)
	ctx.JSON(httpCode, ResponseWithTotalCount{
		RetCode:    retCode,
		Message:    msg,
		DataSet:    data,
		TotalCount: totalCount,
	})
}
