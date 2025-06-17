/*
@Date: 2021/1/12 下午2:31
@Author: max.liu
@File : router
@Desc:
*/

package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/maxliu9403/ProxyHub/internal/common"
	"github.com/maxliu9403/common/middleware"
	"github.com/opentracing/opentracing-go"
)

func RegisterRouter(tra opentracing.Tracer, group *gin.RouterGroup) {
	if tra != nil {
		group.Use(middleware.GinInterceptorWithTrace(tra, true))
	} else {
		group.Use(middleware.GinInterceptor(true))
	}

	proxyGroup := newProxyGroupController(common.BaseController{})
	// 注册分组管理路有
	registerProxyGroupRouter(proxyGroup, group)
}

func registerProxyGroupRouter(proxyGroup *proxyGroupController, group *gin.RouterGroup) {
	group.POST("/api/proxy_group/list", proxyGroup.GetList)
	group.GET("/api/proxy_group/:id", proxyGroup.GetDetail)
	group.DELETE("/api/proxy_group/delete", proxyGroup.Delete)
	group.POST("/api/proxy_group", proxyGroup.Create)
}
