/*
@Date: 2021/1/12 下午2:31
@Author: max.liu
@File : router
@Desc:
*/

package handler

import (
	"fmt"
	"github.com/opentracing/opentracing-go"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/maxliu9403/ProxyHub/internal/common"
	"github.com/maxliu9403/common/middleware"
)

var (
	base = common.BaseController{}

	// controller 以及 handleMap, 新增 action 时修改之
	demoCtrl = newDemoController(base)

	actionHandleFunc = map[string]func(c *gin.Context){
		"GetList":   demoCtrl.GetList,
		"GetDetail": demoCtrl.GetDetail,
		"Delete":    demoCtrl.Delete,
	}
)

type RawAction struct {
	Action string
}

func RegisterRouter(tra opentracing.Tracer, group *gin.RouterGroup) {
	if tra != nil {
		group.Use(middleware.GinInterceptorWithTrace(tra, true))
	} else {
		group.Use(middleware.GinInterceptor(true))
	}

	{
		group.POST("", handler)
	}
}

func handler(c *gin.Context) {
	/* 优先使用query里面的Action，再使用body里面的Action */
	var actionName string
	queryAction := c.Query("Action")

	if queryAction == "" {
		var rawAction RawAction
		_ = c.ShouldBindBodyWith(&rawAction, binding.JSON)
		if rawAction.Action == "" {
			base.Response(c, nil, common.NewErrorCode(common.ActionNotFound, fmt.Errorf("请检查请求")))
			return
		}

		actionName = rawAction.Action
	} else {
		actionName = queryAction
	}

	action, ok := actionHandleFunc[actionName]
	if !ok {
		base.Response(c, nil, common.NewErrorCode(common.ActionNotFound, fmt.Errorf("请检查请求")))
	} else {
		action(c)
	}
}
