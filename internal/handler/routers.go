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

func RegisterRouter(tra opentracing.Tracer, r *gin.RouterGroup) {
	if tra != nil {
		r.Use(middleware.GinInterceptorWithTrace(tra, true))
	} else {
		r.Use(middleware.GinInterceptor(true))
	}

	// 管理员鉴权中间件
	adminMW := middleware.AdminAuthMiddleware(secret)

	// 控制器初始化
	groupCtl := newGroupController(common.BaseController{})
	proxyCtl := newProxyController(common.BaseController{})
	tokenCtl := newTokenController(common.BaseController{})
	emulatorCtl := newEmulatorController(common.BaseController{})

	// 管理员接口路由（带中间件）
	adminGroup := r.Group("/api")
	adminGroup.Use(adminMW)
	registerProxyRouter(proxyCtl, adminGroup)
	registerTokenRouter(tokenCtl, adminGroup)
	registerGroupRouter(groupCtl, adminGroup)
	registerEmulatorRouter(emulatorCtl, adminGroup)
}

func registerGroupRouter(proxyGroup *groupController, group *gin.RouterGroup) {
	group.POST("/group/search", proxyGroup.GetList)
	group.GET("/group/:id", proxyGroup.Detail)
	group.DELETE("/group", proxyGroup.Delete)
	group.POST("/group", proxyGroup.Create)
	group.PUT("/group", proxyGroup.Update)
}

func registerProxyRouter(proxy *proxyController, group *gin.RouterGroup) {
	group.POST("/proxy/search", proxy.GetList)
	group.GET("/proxy/:ip", proxy.Detail)
	group.DELETE("/proxy", proxy.Delete)
	group.POST("/proxy", proxy.Create)
	group.PUT("/proxy", proxy.Update)
}

func registerTokenRouter(token *tokenController, group *gin.RouterGroup) {
	group.POST("/token/search", token.GetList)
	group.DELETE("/token", token.Delete)
	group.POST("/token", token.Create)
	group.GET("/token/validate", token.Validate)
}

func registerEmulatorRouter(emulator *emulatorController, group *gin.RouterGroup) {
	group.POST("/emulator/search", emulator.GetList)
	group.DELETE("/emulator", emulator.Delete)
	group.POST("/emulator", emulator.Create)
	group.GET("/emulator:uuid", emulator.Detail)
}
