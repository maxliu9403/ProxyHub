/*
@Date: 2025/6/17 18:56
@Author: max.liu
@File : proxy_group
*/

package handler

import (
	"fmt"
	"github.com/maxliu9403/ProxyHub/internal/types"
	"strconv"
	"strings"

	"github.com/maxliu9403/ProxyHub/models"

	"github.com/gin-gonic/gin"
	"github.com/maxliu9403/ProxyHub/internal/common"
	"github.com/maxliu9403/ProxyHub/internal/logic/proxy_group"
)

type proxyGroupController struct {
	common.BaseController
}

func newProxyGroupController(base common.BaseController) *proxyGroupController {
	return &proxyGroupController{BaseController: base}
}

// GetList godoc
// @Summary     获取分组列表
// @Description 获取分组列表
// @Tags        分组管理
// @Accept      json
// @Produce     json
// @Param   	params     	body    	types.BasicQuery     	false    "查询通用请求参数"
// @Success     200     {object}        common.ResponseWithTotalCount{Data=[]models.ProxyGroups} "结果：{RetCode:code,Data:数据,Message:消息}"
// @Failure     500     {object}        common.Response "结果：{RetCode:code,Data:数据,Message:消息}"
// @Router      /api/proxy_group/list [post]
func (m *proxyGroupController) GetList(c *gin.Context) {
	var (
		svc    proxy_group.Svc
		err    error
		params types.BasicQuery
	)

	if !m.CheckParams(c, &params) {
		return
	}

	if params.Limit == 0 {
		params.Limit = 10
	}

	params.Keyword = strings.TrimSpace(params.Keyword)

	svc.Ctx = c
	data, err := svc.GetList(params)
	if data == nil || err != nil {
		m.ResponseWithTotalCount(c, []models.ProxyGroups{}, 0, common.NewErrorCode(common.ErrGetList, err))
		return
	}
	m.ResponseWithTotalCount(c, data, 0, common.NewErrorCode(common.ErrGetList, err))
}

// GetDetail godoc
// @Summary     获取分组详情
// @Description 获取分组详情
// @Tags        分组管理
// @Accept      json
// @Produce     json
// @Param       id   path     int  true  "分组ID"
// @Success     200     {object}        common.Response{Data=models.ProxyGroups} "结果：{RetCode:code,Data:数据,Message:消息}"
// @Failure     500     {object}        common.Response "结果：{RetCode:code,Data:数据,Message:消息}"
// @Router      /api/proxy_group/{id} [get]
func (m *proxyGroupController) GetDetail(c *gin.Context) {
	var (
		svc proxy_group.Svc
		err error
	)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		m.Response(c, nil, common.NewErrorCode(common.ErrInvalidParams, fmt.Errorf("无效的ID参数")))
		return
	}

	svc.Ctx = c
	svc.ID = id
	resp, err := svc.Detail()
	m.Response(c, resp, common.NewErrorCode(common.ErrGetDetail, err))
}

// Delete godoc
// @Summary     删除分组
// @Description 删除分组
// @Tags        分组管理
// @Accept      json
// @Produce     json
// @Param   	params     	body    	proxy_group.DeleteParams     	false    "删除请求参数"
// @Success     200     {object}        common.Response "结果：{RetCode:code,Data:数据,Message:消息}"
// @Failure     500     {object}        common.Response "结果：{RetCode:code,Data:数据,Message:消息}"
// @Router      /api/proxy_group/delete [delete]
func (m *proxyGroupController) Delete(c *gin.Context) {
	var (
		svc    proxy_group.Svc
		err    error
		params proxy_group.DeleteParams
	)

	if !m.CheckParams(c, &params) {
		return
	}

	svc.Ctx = c
	svc.RunningTest = params.Test.Enable
	err = svc.Delete(params)
	m.Response(c, nil, common.NewErrorCode(common.ErrDelete, err))
}

// Create godoc
// @Summary     创建分组
// @Description 创建一个新的代理分组
// @Tags        分组管理
// @Accept      json
// @Produce     json
// @Param       params  body  proxy_group.CreateParams  true  "创建参数"
// @Success     200     {object}  common.Response{Data=models.ProxyGroups}  "结果：{RetCode:code,Data:数据,Message:消息}"
// @Failure     500     {object}  common.Response
// @Router      /api/proxy_group [post]
func (m *proxyGroupController) Create(c *gin.Context) {
	var (
		svc    proxy_group.Svc
		err    error
		params proxy_group.CreateParams
	)

	if !m.CheckParams(c, &params) {
		return
	}

	svc.Ctx = c
	svc.RunningTest = params.Test.Enable

	resp, err := svc.Create(params)
	m.Response(c, resp, common.NewErrorCode(common.ErrCreate, err))
}
