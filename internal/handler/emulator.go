package handler

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/maxliu9403/ProxyHub/internal/common"
	"github.com/maxliu9403/ProxyHub/internal/logic/emulator"
	"github.com/maxliu9403/ProxyHub/internal/types"
	"github.com/maxliu9403/ProxyHub/models"
)

type emulatorController struct {
	common.BaseController
}

func newEmulatorController(base common.BaseController) *emulatorController {
	return &emulatorController{BaseController: base}
}

// GetList godoc
// @Summary     获取模拟器列表
// @Description 获取模拟器列表
// @Tags        模拟器管理
// @Security    AdminTokenAuth
// @Accept      json
// @Produce     json
// @Param       params body types.BasicQuery false "查询通用请求参数"
// @Success     200 {object} common.ResponseWithTotalCount{Data=[]models.Emulator}
// @Failure     500 {object} common.Response
// @Router      /api/emulator/search [post]
func (e *emulatorController) GetList(c *gin.Context) {
	var (
		svc    emulator.Svc
		params types.BasicQuery
	)

	if !e.CheckParams(c, &params) {
		return
	}

	params.Keyword = strings.TrimSpace(params.Keyword)
	svc.Ctx = c

	data, err := svc.GetList(params)
	if err != nil || data == nil {
		e.ResponseWithTotalCount(c, []models.Emulator{}, 0, err)
		return
	}

	e.ResponseWithTotalCount(c, data.Data, data.Counts, nil)
}

// Create godoc
// @Summary     批量创建模拟器
// @Description 批量创建模拟器（自动跳过已存在）
// @Tags        模拟器管理
// @Security    AdminTokenAuth
// @Accept      json
// @Produce     json
// @Param       params body emulator.CreateBatchParams true "批量创建参数"
// @Success     200 {object} common.Response{Data=emulator.CreateBatchResp}
// @Failure     500 {object} common.Response
// @Router      /api/emulator [post]
func (e *emulatorController) Create(c *gin.Context) {
	var (
		svc    emulator.Svc
		params emulator.CreateBatchParams
	)

	if !e.CheckParams(c, &params) {
		return
	}

	svc.Ctx = c
	resp, err := svc.CreateBatch(params)
	e.Response(c, resp, err)
}

// Detail godoc
// @Summary     获取模拟器详情
// @Description 根据 UUID 获取模拟器详情
// @Tags        模拟器管理
// @Security    AdminTokenAuth
// @Accept      json
// @Produce     json
// @Param       uuid query string true "UUID"
// @Success     200 {object} common.Response{Data=models.Emulator}
// @Failure     500 {object} common.Response
// @Router      /api/emulator/{uuid} [get]
func (e *emulatorController) Detail(c *gin.Context) {
	uuid := c.Query("uuid")
	if uuid == "" {
		e.Response(c, nil, common.NewErrorCode(common.ErrInvalidParams, fmt.Errorf("无效的UUID参数")))
		return
	}

	var svc emulator.Svc
	svc.Ctx = c

	data, err := svc.Detail(uuid)
	e.Response(c, data, err)
}

// Delete godoc
// @Summary     删除模拟器
// @Description 批量删除模拟器
// @Tags        模拟器管理
// @Security    AdminTokenAuth
// @Accept      json
// @Produce     json
// @Param       params body emulator.DeleteParams true "删除参数"
// @Success     200 {object} common.Response
// @Failure     500 {object} common.Response
// @Router      /api/emulator [delete]
func (e *emulatorController) Delete(c *gin.Context) {
	var (
		svc    emulator.Svc
		params emulator.DeleteParams
	)

	if !e.CheckParams(c, &params) {
		return
	}

	svc.Ctx = c
	err := svc.Delete(params)
	e.Response(c, nil, err)
}

// Update godoc
// @Summary     更新模拟器分组
// @Description 更新模拟器分组信息
// @Tags        模拟器管理
// @Security    AdminTokenAuth
// @Accept      json
// @Produce     json
// @Param       params body emulator.UpdateParams true "更新参数"
// @Success     200 {object} common.Response
// @Failure     400 {object} common.Response "参数错误"
// @Failure     500 {object} common.Response "服务器内部错误"
// @Router      /api/emulator [put]
func (e *emulatorController) Update(c *gin.Context) {
	var (
		svc    emulator.Svc
		params emulator.UpdateParams
	)

	if !e.CheckParams(c, &params) {
		return
	}

	svc.Ctx = c
	err := svc.Update(params)
	e.Response(c, nil, err)
}
