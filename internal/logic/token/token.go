package token

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/maxliu9403/ProxyHub/internal/common"
	"github.com/maxliu9403/ProxyHub/internal/logic"
	"github.com/maxliu9403/ProxyHub/internal/types"
	"github.com/maxliu9403/ProxyHub/models"
	"github.com/maxliu9403/ProxyHub/models/factory"
	"github.com/maxliu9403/ProxyHub/models/repo"
	"github.com/maxliu9403/common/gormdb"
	"github.com/maxliu9403/common/logger"
	"gorm.io/gorm"
)

type Svc struct {
	ID          int64
	Ctx         context.Context
	RunningTest bool
	DB          *gorm.DB
}

func (s *Svc) getRepo() repo.TokenRepo {
	s.DB = gormdb.Cli(s.Ctx)
	return factory.TokenRepo(s.DB)
}

type CreateParams struct {
	Description string `json:"Description" binding:"required"` // 描述
	Expired     *int64 `json:"Expired,omitempty"`              // 到期时间，不传就是永久不过期
}

func (p CreateParams) ToModel(token string) *models.Token {
	var expireAt *time.Time
	if p.Expired != nil && *p.Expired > 0 {
		t := time.Unix(*p.Expired, 0)
		expireAt = &t
	}

	return &models.Token{
		Token:       token,
		Description: p.Description,
		ExpireAt:    expireAt,
	}
}

func (s *Svc) Create(params CreateParams) (*models.Token, error) {
	// 校验过期时间
	if params.Expired != nil {
		now := time.Now().Unix()
		if *params.Expired < now {
			return nil, common.NewErrorCode(common.ErrBuildToken, fmt.Errorf("过期时间不能小于当前时间"))
		}
	}

	// 实现Token的生成
	token, err := logic.GenerateSecureToken(32)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, "generate token failed: %s", err.Error())
		return nil, common.NewErrorCode(common.ErrBuildToken, err)
	}

	model := params.ToModel(token)
	err = s.getRepo().Create(model)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, "create failed: %s", err.Error())
		return nil, common.NewErrorCode(common.ErrCreateGroup, err)
	}

	return model, nil
}

type DeleteToken struct {
	common.Test
	Tokens []string `json:"Tokens" binding:"required"` // 待删除 Token 列表
}

func (s *Svc) Delete(params DeleteToken) (err error) {
	crud := s.getRepo()
	err = crud.Deletes(params.Tokens)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, "delete %d failed: %s", params.Tokens, err.Error())
		return common.NewErrorCode(common.ErrDeleteGroup, err)
	}

	return err
}

func (s *Svc) GetList(q types.BasicQuery) (data *common.ListData, err error) {
	data = &common.ListData{}

	crud := s.getRepo()
	table := &models.Token{}
	tokenList := make([]models.Token, 0)
	total, err := crud.GetList(q, table, &tokenList)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, "query list failed: %s", err.Error())
		return data, common.NewErrorCode(common.ErrGetList, err)
	}

	data.Counts = total
	data.Data = tokenList

	return data, err
}

func (s *Svc) ValidateToken(token string) (bool, error) {
	isValid, err := s.getRepo().IsValid(token)
	if err != nil {
		return false, err
	}
	if isValid {
		return true, nil
	}
	return false, errors.New("token已过期")
}
