/*
@Date: 2025/6/17
@Author: max.liu
@File : token
*/

package handler

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/maxliu9403/ProxyHub/internal/common"
	"github.com/maxliu9403/ProxyHub/internal/logic/token"
	"github.com/maxliu9403/ProxyHub/models"
)

type tokenController struct {
	common.BaseController
}

func newTokenController(base common.BaseController) *tokenController {
	return &tokenController{BaseController: base}
}

// Create godoc
// @Summary     创建 Token
// @Description 创建一个新的访问 Token，如果当前分组内存在未过期的Token则会删除历史Token（需管理员权限）
// @Tags        Token 管理
// @Security    AdminTokenAuth
// @Accept      json
// @Produce     json
// @Param       params  body  token.CreateParams  true  "创建参数"
// @Success     200     {object}  common.Response{Data=models.Token}
// @Failure     500     {object}  common.Response
// @Router      /api/token [post]
func (m *tokenController) Create(c *gin.Context) {
	var (
		svc    token.Svc
		err    error
		params token.CreateParams
	)

	if !m.CheckParams(c, &params) {
		return
	}

	svc.Ctx = c
	resp, err := svc.Create(params)
	m.Response(c, resp, common.NewErrorCode(common.ErrCreateToken, err))
}

// Delete godoc
// @Summary     删除 Token
// @Description 删除指定 ID 的 Token（需管理员权限）
// @Tags        Token 管理
// @Security    AdminTokenAuth
// @Accept      json
// @Produce     json
// @Param       params  body  token.DeleteToken  true  "删除参数"
// @Success     200 {object}  common.Response
// @Failure     500 {object}  common.Response
// @Router      /api/token [delete]
func (m *tokenController) Delete(c *gin.Context) {
	var (
		svc    token.Svc
		err    error
		params token.DeleteToken
	)

	if !m.CheckParams(c, &params) {
		return
	}

	svc.Ctx = c
	err = svc.Delete(params)
	m.Response(c, nil, common.NewErrorCode(common.ErrDeleteToken, err))
}

// GetList godoc
// @Summary     获取 Token 列表
// @Description 获取所有 Token（需管理员权限）
// @Tags        Token 管理
// @Security    AdminTokenAuth
// @Accept      json
// @Produce     json
// @Param   	params     	body    	models.GetTokenListParams     	false    "查询通用请求参数"
// @Success     200 {object}  common.Response{Data=[]models.Token}
// @Failure     500 {object}  common.Response
// @Router      /api/token/search [post]
func (m *tokenController) GetList(c *gin.Context) {
	var (
		svc    token.Svc
		err    error
		params models.GetTokenListParams
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
		m.ResponseWithTotalCount(c, []models.Token{}, 0, common.NewErrorCode(common.ErrGetList, err))
		return
	}
	m.ResponseWithTotalCount(c, data.Data, data.Counts, common.NewErrorCode(common.ErrGetList, err))
}

// Validate godoc
// @Summary     校验 Token
// @Description 校验 Token 是否有效（用于订阅接口）
// @Tags        Token 管理
// @Security    AdminTokenAuth
// @Accept      json
// @Produce     json
// @Param       token  query  string  true  "待校验的 Token"
// @Success     200 {object}  common.Response{Data=bool} "结果 true 表示有效"
// @Failure     500 {object}  common.Response
// @Router      /api/token/validate [get]
func (m *tokenController) Validate(c *gin.Context) {
	var (
		svc token.Svc
		err error
	)

	q := c.Query("token")
	if len(q) == 0 {
		m.Response(c, nil, common.NewErrorCode(common.ErrInvalidParams, fmt.Errorf("无效的Token")))
		return
	}

	svc.Ctx = c
	ok, err := svc.ValidateToken(q)
	m.Response(c, ok, common.NewErrorCode(common.ErrValidateToken, err))
}
