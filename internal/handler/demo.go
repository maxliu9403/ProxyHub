/*
@Date: 2021/12/8 18:56
@Author: max.liu
@File : monitor
*/

package handler

import (
	"github.com/maxliu9403/go-template/internal/types"
	"strings"

	"github.com/maxliu9403/go-template/models"

	"github.com/gin-gonic/gin"
	"github.com/maxliu9403/go-template/internal/common"
	"github.com/maxliu9403/go-template/internal/logic/demo"
)

type demoController struct {
	common.BaseController
}

func newDemoController(base common.BaseController) *demoController {
	return &demoController{BaseController: base}
}

// GetList godoc
// @Summary     获取列表
// @Description 获取列表
// @Tags        Demo
// @Accept      json
// @Produce     json
// @Param   	params     	body    	types.BasicQuery     	false    "查询通用请求参数"
// @Success     200     {object}        common.ResponseWithTotalCount{Data=[]models.Demo} "结果：{RetCode:code,Data:数据,Message:消息}"
// @Failure     500     {object}        common.Response "结果：{RetCode:code,Data:数据,Message:消息}"
// @Router      ?Action=GetList             [post]
func (m *demoController) GetList(c *gin.Context) {
	var (
		svc    demo.Svc
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
		m.ResponseWithTotalCount(c, []models.Demo{}, 0, common.NewErrorCode(common.ErrGetList, err))
		return
	}
	m.ResponseWithTotalCount(c, data, 0, common.NewErrorCode(common.ErrGetList, err))
}

// GetDetail godoc
// @Summary     获取详情
// @Description 获取详情
// @Tags        Demo
// @Accept      json
// @Produce     json
// @Param   	params     	body    	demo.IDParams     	false    "请求参数"
// @Success     200     {object}        common.Response{Data=models.Demo} "结果：{RetCode:code,Data:数据,Message:消息}"
// @Failure     500     {object}        common.Response "结果：{RetCode:code,Data:数据,Message:消息}"
// @Router      ?Action=GetDetail             [post]
func (m *demoController) GetDetail(c *gin.Context) {
	var (
		svc    demo.Svc
		err    error
		params demo.IDParams
	)

	if !m.CheckParams(c, &params) {
		return
	}

	svc.Ctx = c
	svc.ID = params.ID
	svc.RunningTest = params.Test.Enable
	resp, err := svc.Detail()
	m.Response(c, resp, common.NewErrorCode(common.ErrGetDetail, err))
}

// Delete Doc
// @Summary     删除 agent
// @Description 删除 agent
// @Tags        Demo
// @Accept      json
// @Produce     json
// @Param   	params     	body    	demo.DeleteParams     	false    "删除请求参数"
// @Success     200     {object}        common.Response "结果：{RetCode:code,Data:数据,Message:消息}"
// @Failure     500     {object}        common.Response "结果：{RetCode:code,Data:数据,Message:消息}"
// @Router      ?Action=Delete             [post]
func (m *demoController) Delete(c *gin.Context) {
	var (
		svc    demo.Svc
		err    error
		params demo.DeleteParams
	)

	if !m.CheckParams(c, &params) {
		return
	}

	svc.Ctx = c
	svc.RunningTest = params.Test.Enable
	err = svc.Delete(params)
	m.Response(c, nil, common.NewErrorCode(common.ErrDelete, err))
}
