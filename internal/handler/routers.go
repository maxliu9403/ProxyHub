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

	groupCtl := newGroupController(common.BaseController{})
	// 注册分组管理路有
	registerGroupRouter(groupCtl, group)

	proxyCtl := newProxyController(common.BaseController{})
	// 注册代理管理路有
	registerProxyRouter(proxyCtl, group)
}

func registerGroupRouter(proxyGroup *groupController, group *gin.RouterGroup) {
	group.POST("/api/group/list", proxyGroup.GetList)
	group.GET("/api/group/:id", proxyGroup.GetDetail)
	group.DELETE("/api/group/delete", proxyGroup.Delete)
	group.POST("/api/group", proxyGroup.Create)
	group.PUT("/api/group", proxyGroup.Update)
}

func registerProxyRouter(proxyGroup *proxyController, group *gin.RouterGroup) {
	group.POST("/api/proxy/list", proxyGroup.GetList)
	group.GET("/api/proxy/:id", proxyGroup.GetDetail)
	group.DELETE("/api/proxy/delete", proxyGroup.Delete)
	group.POST("/api/proxy", proxyGroup.Create)
	group.PUT("/api/proxy", proxyGroup.Update)
}
