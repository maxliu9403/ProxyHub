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

const secret = "x-secret"

func RegisterRouter(tra opentracing.Tracer, group *gin.RouterGroup) {
	if tra != nil {
		group.Use(middleware.GinInterceptorWithTrace(tra, true))
	} else {
		group.Use(middleware.GinInterceptor(true))
	}

	// 管理员鉴权中间件
	adminMW := middleware.AdminAuthMiddleware(secret)

	// 控制器初始化
	groupCtl := newGroupController(common.BaseController{})
	proxyCtl := newProxyController(common.BaseController{})
	tokenCtl := newTokenController(common.BaseController{})

	// 管理员接口路由（带中间件）
	adminGroup := group.Group("")
	adminGroup.Use(adminMW)
	registerProxyRouter(proxyCtl, adminGroup)
	registerTokenRouter(tokenCtl, adminGroup)
	registerGroupRouter(groupCtl, adminGroup)
}

func registerGroupRouter(proxyGroup *groupController, group *gin.RouterGroup) {
	group.POST("/api/group/list", proxyGroup.GetList)
	group.GET("/api/group/:id", proxyGroup.GetDetail)
	group.DELETE("/api/group/delete", proxyGroup.Delete)
	group.POST("/api/group", proxyGroup.Create)
	group.PUT("/api/group", proxyGroup.Update)
}

func registerProxyRouter(proxy *proxyController, group *gin.RouterGroup) {
	group.POST("/api/proxy/list", proxy.GetList)
	group.POST("/api/proxy/detail", proxy.GetDetail)
	group.DELETE("/api/proxy/delete", proxy.Delete)
	group.POST("/api/proxy", proxy.Create)
	group.PUT("/api/proxy", proxy.Update)
}

func registerTokenRouter(token *tokenController, group *gin.RouterGroup) {
	group.POST("/api/token/list", token.GetList)
	group.DELETE("/api/token/delete", token.Delete)
	group.POST("/api/token", token.Create)
	group.GET("/api/token/validate", token.Validate)
}
