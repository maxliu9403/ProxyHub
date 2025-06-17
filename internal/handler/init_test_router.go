/*
@Date: 2022/4/22 12:28
@Author: max.liu
@File : init_test_router
@Desc:
*/

package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/maxliu9403/ProxyHub/internal/config"
	"github.com/maxliu9403/common/apiserver"
)

func setupRouter() *gin.Engine {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config.G.App.RunMode = apiserver.RunModeRelease
	config.G.Log.Level = "debug"

	s := apiserver.CreateNewServer(ctx, config.G.APIConfig)
	defer s.Stop()

	v1Group := s.AddGinGroup("")
	RegisterRouter(nil, v1Group)

	return s.ExposeEng()
}
