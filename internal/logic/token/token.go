package token

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/maxliu9403/ProxyHub/internal/logic"
	"github.com/maxliu9403/ProxyHub/internal/logic/group"

	"github.com/maxliu9403/ProxyHub/internal/common"
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

type Validator interface {
	ValidateToken(token string) (bool, error)
}

type CreateParams struct {
	GroupID     int64  `json:"GroupID" binding:"required,gt=0" comment:"所属分组ID"`
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
		GroupID:     p.GroupID,
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

	// 校验分组
	groupAPI := group.NewGroupAPI(s.Ctx)
	hasGroup, err := groupAPI.CheckGroupID(params.GroupID)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, "group check failed: %s", err.Error())
		return nil, common.NewErrorCode(common.ErrBuildTokenErrGroup, err)
	}
	if !hasGroup {
		return nil, errors.New("当前分组ID不是激活状态")
	}

	// 启动事务
	tx := gormdb.Cli(s.Ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	tokenRepo := factory.TokenRepo(tx)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 删除旧 Token
	now := time.Now()
	oldTokens, err := tokenRepo.GetValidTokensByGroup(params.GroupID, now)
	if err != nil {
		tx.Rollback()
		return nil, common.NewErrorCode(common.ErrGetList, err)
	}
	if len(oldTokens) > 0 {
		tokens := make([]string, 0, len(oldTokens))
		for _, t := range oldTokens {
			tokens = append(tokens, t.Token)
		}
		if err := tokenRepo.Deletes(tokens); err != nil {
			tx.Rollback()
			return nil, common.NewErrorCode(common.ErrDeleteGroup, err)
		}
	}

	// 创建新 Token
	tokenStr, err := logic.GenerateSecureToken(32)
	if err != nil {
		tx.Rollback()
		return nil, common.NewErrorCode(common.ErrBuildToken, err)
	}

	model := params.ToModel(tokenStr)
	if err := tokenRepo.Create(model); err != nil {
		tx.Rollback()
		return nil, common.NewErrorCode(common.ErrCreateGroup, err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
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

type GetListParams struct {
	types.BasicQuery         // Limit, Offset, Keyword, Order 等
	GroupIDs         []int64 `json:"GroupIDs,omitempty"` // 多组 ID 过滤
}

func (s *Svc) GetList(q models.GetTokenListParams) (data *common.ListData, err error) {
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
