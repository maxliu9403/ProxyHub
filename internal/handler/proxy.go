/*
@Date: 2025/6/17
@Author: max.liu
@File : proxy
*/

package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/maxliu9403/ProxyHub/internal/common"
	"github.com/maxliu9403/ProxyHub/internal/logic/proxy"
	"github.com/maxliu9403/ProxyHub/models"
)

type proxyController struct {
	common.BaseController
}

func newProxyController(base common.BaseController) *proxyController {
	return &proxyController{BaseController: base}
}

// GetList godoc
// @Summary     获取代理列表
// @Description 支持分页与多条件查询
// @Tags        代理管理
// @Security    AdminTokenAuth
// @Accept      json
// @Produce     json
// @Param       params body models.GetListParams false "查询参数"
// @Success     200 {object} common.ResponseWithTotalCount{Data=[]models.Proxy} "结果：{RetCode:code,Data:数据,Message:消息}"
// @Failure     500 {object} common.Response
// @Router      /api/proxy/search [post]
func (m *proxyController) GetList(c *gin.Context) {
	var (
		svc    proxy.Svc
		err    error
		params models.GetListParams
	)

	if !m.CheckParams(c, &params) {
		return
	}
	if params.Limit == 0 {
		params.Limit = 10
	}

	svc.Ctx = c
	resp, err := svc.GetList(params)
	if err != nil || resp == nil {
		m.ResponseWithTotalCount(c, []models.Proxy{}, 0, common.NewErrorCode(common.ErrGetList, err))
		return
	}
	m.ResponseWithTotalCount(c, resp.Data, resp.Counts, nil)
}

// Detail godoc
// @Summary     获取代理详情
// @Description 通过 IP 获取代理信息（单个）
// @Tags        代理管理
// @Security    AdminTokenAuth
// @Accept      json
// @Produce     json
// @Param       ip   path     string  true  "代理 IP"
// @Success     200  {object}  common.Response{Data=models.Proxy}
// @Failure     500  {object}  common.Response
// @Router      /api/proxy/{ip} [get]
func (m *proxyController) Detail(c *gin.Context) {
	ip := c.Param("ip")
	if ip == "" {
		m.Response(c, nil, common.NewErrorCode(common.ErrInvalidParams, fmt.Errorf("无效的IP参数")))
		return
	}

	svc := proxy.Svc{Ctx: c}
	proxyItem, err := svc.GetByIP(ip)
	m.Response(c, proxyItem, common.NewErrorCode(common.ErrGetDetail, err))
}

// Create godoc
// @Summary     批量创建代理
// @Description 创建一个或多个新的代理IP记录
// @Tags        代理管理
// @Security    AdminTokenAuth
// @Accept      json
// @Produce     json
// @Param       params  body  proxy.CreateBatchParams  true  "创建参数数组"
// @Success     200     {object}  common.Response{Data=[]models.Proxy}
// @Failure     500     {object}  common.Response
// @Router      /api/proxy [post]
func (m *proxyController) Create(c *gin.Context) {
	var (
		svc    proxy.Svc
		err    error
		params proxy.CreateBatchParams
	)

	if !m.CheckParams(c, &params) {
		return
	}

	svc.Ctx = c

	resp, err := svc.CreateBatch(params)
	m.Response(c, resp, common.NewErrorCode(common.ErrCreateProxy, err))
}

// Update godoc
// @Summary     更新代理
// @Description 更新代理信息
// @Tags        代理管理
// @Security    AdminTokenAuth
// @Accept      json
// @Produce     json
// @Param       params  body  proxy.UpdateParams  true  "更新参数"
// @Success     200     {object}  common.Response{Data=models.Proxy}
// @Failure     500     {object}  common.Response
// @Router      /api/proxy [put]
func (m *proxyController) Update(c *gin.Context) {
	var (
		svc    proxy.Svc
		err    error
		params proxy.UpdateParams
	)

	if !m.CheckParams(c, &params) {
		return
	}

	svc.Ctx = c

	err = svc.Update(params)
	m.Response(c, nil, common.NewErrorCode(common.ErrUpdateProxy, err))
}

// Delete godoc
// @Summary     删除代理
// @Description 删除一个或多个代理
// @Tags        代理管理
// @Security    AdminTokenAuth
// @Accept      json
// @Produce     json
// @Param       params  body  proxy.DeleteParams  true  "删除请求参数"
// @Success     200     {object}  common.Response
// @Failure     500     {object}  common.Response
// @Router      /api/proxy [delete]
func (m *proxyController) Delete(c *gin.Context) {
	var (
		svc    proxy.Svc
		err    error
		params proxy.DeleteParams
	)

	if !m.CheckParams(c, &params) {
		return
	}

	svc.Ctx = c
	err = svc.Delete(params)
	m.Response(c, nil, common.NewErrorCode(common.ErrDeleteProxy, err))
}
