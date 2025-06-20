package handler

import (
	"fmt"
	"net/http"

	"github.com/maxliu9403/ProxyHub/internal/logic/token"

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
// @Description 通过 token 和 uuid 获取对应的 Clash 代理配置（YAML 格式）
// @Tags        订阅管理
// @Produce     plain
// @Param       token   path     string  true  "授权 Token"
// @Param       uuid    path     string  true  "模拟器 uuid"
// @Success     200     {string}  string  "YAML 配置内容"
// @Failure     400     {object} common.Response "参数错误"
// @Failure     500     {object} common.Response "服务器内部错误"
// @Router      /api/subscribe/{token}/{uuid} [get]
func (m *subscribeController) Get(c *gin.Context) {
	tokenParam := c.Param("token")
	uuid := c.Param("uuid")

	if tokenParam == "" || uuid == "" {
		m.Response(c, nil, common.NewErrorCode(common.ErrInvalidParams, fmt.Errorf("存在无效参数")))
		return
	}

	tokenSvc := &token.Svc{Ctx: c} // 需要是指针，因为接口是由 *token.Svc 实现的
	svc := subscribe.Svc{
		Ctx:            c,
		TokenValidator: tokenSvc,
	}
	clashCfg, err := svc.Subscribe(tokenParam, uuid)
	if err != nil {
		m.Response(c, nil, common.NewErrorCode(common.ErrGetSubscribe, err))
		return
	}

	// 设置响应头并直接写入 YAML 配置内容
	c.Header("Content-Type", "application/yaml")
	c.String(http.StatusOK, clashCfg)
}
