/*
@Date: 2025/6/17 18:56
@Author: max.liu
@File : group
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
	"github.com/maxliu9403/ProxyHub/internal/logic/group"
)

type groupController struct {
	common.BaseController
}

func newGroupController(base common.BaseController) *groupController {
	return &groupController{BaseController: base}
}

// GetList godoc
// @Summary     获取分组列表
// @Description 获取分组列表
// @Tags        分组管理
// @Accept      json
// @Produce     json
// @Param   	params     	body    	types.BasicQuery     	false    "查询通用请求参数"
// @Success     200     {object}        common.ResponseWithTotalCount{Data=[]models.Groups} "结果：{RetCode:code,Data:数据,Message:消息}"
// @Failure     500     {object}        common.Response "结果：{RetCode:code,Data:数据,Message:消息}"
// @Router      /api/group/list [post]
func (m *groupController) GetList(c *gin.Context) {
	var (
		svc    group.Svc
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
		m.ResponseWithTotalCount(c, []models.Groups{}, 0, common.NewErrorCode(common.ErrGetList, err))
		return
	}
	m.ResponseWithTotalCount(c, data.Data, data.Counts, common.NewErrorCode(common.ErrGetList, err))
}

// GetDetail godoc
// @Summary     获取分组详情
// @Description 获取分组详情
// @Tags        分组管理
// @Accept      json
// @Produce     json
// @Param       id   path     int  true  "分组ID"
// @Success     200     {object}        common.Response{Data=models.Groups} "结果：{RetCode:code,Data:数据,Message:消息}"
// @Failure     500     {object}        common.Response "结果：{RetCode:code,Data:数据,Message:消息}"
// @Router      /api/group/{id} [get]
func (m *groupController) GetDetail(c *gin.Context) {
	var (
		svc group.Svc
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
	fmt.Println(svc.ID, "===")
	resp, err := svc.Detail()
	m.Response(c, resp, common.NewErrorCode(common.ErrGetDetail, err))
}

// Delete godoc
// @Summary     删除分组
// @Description 删除分组
// @Tags        分组管理
// @Accept      json
// @Produce     json
// @Param   	params     	body    	group.DeleteParams     	false    "删除请求参数"
// @Success     200     {object}        common.Response "结果：{RetCode:code,Data:数据,Message:消息}"
// @Failure     500     {object}        common.Response "结果：{RetCode:code,Data:数据,Message:消息}"
// @Router      /api/group/delete [delete]
func (m *groupController) Delete(c *gin.Context) {
	var (
		svc    group.Svc
		err    error
		params group.DeleteParams
	)

	if !m.CheckParams(c, &params) {
		return
	}

	svc.Ctx = c
	svc.RunningTest = params.Test.Enable
	err = svc.Delete(params)
	m.Response(c, nil, common.NewErrorCode(common.ErrDeleteGroup, err))
}

// Create godoc
// @Summary     创建分组
// @Description 创建一个新的代理分组
// @Tags        分组管理
// @Accept      json
// @Produce     json
// @Param       params  body  group.CreateParams  true  "创建参数"
// @Success     200     {object}  common.Response{Data=models.Groups}  "结果：{RetCode:code,Data:数据,Message:消息}"
// @Failure     500     {object}  common.Response
// @Router      /api/group [post]
func (m *groupController) Create(c *gin.Context) {
	var (
		svc    group.Svc
		err    error
		params group.CreateParams
	)

	if !m.CheckParams(c, &params) {
		return
	}

	svc.Ctx = c
	svc.RunningTest = params.Test.Enable

	resp, err := svc.Create(params)
	m.Response(c, resp, common.NewErrorCode(common.ErrCreateGroup, err))
}

// Update godoc
// @Summary     更新分组
// @Description 更新一个已有的代理分组
// @Tags        分组管理
// @Accept      json
// @Produce     json
// @Param       params  body  group.UpdateParams  true  "更新参数"
// @Success     200     {object}  common.Response{Data=models.Groups}  "结果：{RetCode:code,Data:数据,Message:消息}"
// @Failure     500     {object}  common.Response
// @Router      /api/group [put]
func (m *groupController) Update(c *gin.Context) {
	var (
		svc    group.Svc
		err    error
		params group.UpdateParams
	)

	if !m.CheckParams(c, &params) {
		return
	}

	svc.Ctx = c
	svc.RunningTest = params.Test.Enable

	err = svc.Update(params)
	m.Response(c, nil, common.NewErrorCode(common.ErrUpdateGroup, err))
}

// TODO 实现查询激活状态的Group
