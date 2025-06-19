package handler

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/maxliu9403/ProxyHub/internal/common"
	"github.com/maxliu9403/ProxyHub/internal/logic/subscribe"
)

type subscribeController struct {
	common.BaseController
}

func newSubscribeController(base common.BaseController) *subscribeController {
	return &subscribeController{BaseController: base}
}

// Get godoc
// @Summary     获取代理配置
// @Description 获取代理配置
// @Tags        代理管理
// @Accept      json
// @Produce     json
// @Param       token   path     int  true  "授权token"
// @Success     200     {object}  common.Response{Data=models.Proxy}
// @Failure     500     {object}  common.Response
// @Router      /api/subscribe/ [get]
func (m *subscribeController) Get(c *gin.Context) {
	var (
		svc subscribe.Svc
		err error
	)
	token := c.Param("token")
	groupIdStr := c.Param("group_id")
	groupId, err := strconv.ParseInt(groupIdStr, 10, 64)
	if err != nil || groupId <= 0 {
		m.Response(c, nil, common.NewErrorCode(common.ErrInvalidParams, fmt.Errorf("无效的group_id参数")))
		return
	}

	svc.Ctx = c
	err = svc.Subscribe(token, groupId)
	//err = svc.Update(params)
	m.Response(c, nil, common.NewErrorCode(common.ErrUpdateProxy, err))
}
