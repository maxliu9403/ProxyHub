package proxy

import (
	"context"
	"errors"
	"fmt"

	"github.com/maxliu9403/ProxyHub/internal/common"
	"github.com/maxliu9403/ProxyHub/internal/logic/group"
	"github.com/maxliu9403/ProxyHub/models"
	"github.com/maxliu9403/ProxyHub/models/factory"
	"github.com/maxliu9403/ProxyHub/models/repo"
	"github.com/maxliu9403/common/gormdb"
	"github.com/maxliu9403/common/logger"
	"gorm.io/gorm"
)

type Svc struct {
	ID  int64
	Ctx context.Context
	DB  *gorm.DB
}

func (s *Svc) getRepo() repo.ProxyRepo {
	s.DB = gormdb.Cli(s.Ctx)
	return factory.ProxyRepo(s.DB)
}

func (s *Svc) getGroupRepo() repo.GroupsRepo {
	s.DB = gormdb.Cli(s.Ctx)
	return factory.GroupsRepo(s.DB)
}

type CreateParams struct {
	IP        string `json:"IP" binding:"required"`
	Port      int64  `json:"Port,omitempty" binding:"omitempty,gt=0,lte=65535"`
	Username  string `json:"Username" binding:"required"`
	Password  string `json:"Password" binding:"required"`
	Source    string `json:"Source" binding:"required"`
	ProxyType string `json:"ProxyType,omitempty" binding:"omitempty,oneof=socks5"`
}

type CreateBatchParams struct {
	GroupID int64          `json:"GroupID" binding:"required,gt=0" comment:"所属分组ID"`
	Proxies []CreateParams `json:"Proxies" binding:"required,dive"` // CreateParams 就是单个 proxy 的结构
}

func (p CreateParams) ToModel(groupID int64) *models.Proxy {
	return &models.Proxy{
		IP:        p.IP,
		Port:      p.Port,
		Username:  p.Username,
		ProxyType: p.ProxyType,
		Password:  p.Password,
		Source:    p.Source,
		GroupID:   groupID,
	}
}

type Invalid struct {
	IP      int64  `json:"IP"`
	Message string `json:"Message"` // 错误信息
}

type CreateBatchResult struct {
	CreatedCount   int       `json:"CreatedCount" comment:"成功创建数量"`
	InvalidProxies []Invalid `json:"InvalidProxies" comment:"无效代理列表"`
}

func (s *Svc) CreateBatch(params CreateBatchParams) (resp *CreateBatchResult, err error) {
	groupAPI := group.NewGroupAPI(s.Ctx)
	hasGroup, err := groupAPI.CheckGroupID(params.GroupID)
	if err != nil {
		return nil, common.NewErrorCode(common.ErrCreateProxyCheckGroup, err)
	}

	if !hasGroup {
		return nil, common.NewErrorCode(common.ErrCreateProxyNotGroup, fmt.Errorf("当前分组ID不是激活状态，或者是不存在的激活ID"))
	}

	invalidProxies := make([]Invalid, 0)
	validProxies := make([]*models.Proxy, 0)
	// 用于存储转换后的模型
	for _, p := range params.Proxies {
		model := p.ToModel(params.GroupID)
		// TODO 实现并发校验ip的有效性
		validProxies = append(validProxies, model)
	}
	// 如果一个合法的都没有
	if len(validProxies) == 0 {
		return nil, common.NewErrorCode(common.ErrCreateProxy, fmt.Errorf("无任何有效代理"))
	}

	// 执行批量创建
	err = s.getRepo().CreateBatch(validProxies)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, "batch create proxy failed: %s", err.Error())
		return nil, common.NewErrorCode(common.ErrCreateProxy, err)
	}
	toCreate := len(params.Proxies) - len(invalidProxies)

	return &CreateBatchResult{
		CreatedCount:   toCreate,
		InvalidProxies: invalidProxies,
	}, nil
}

type UpdateParams struct {
	ID        int64   `json:"ID" binding:"required"`                             // ID 必填
	IP        *string `json:"IP" binding:"omitempty"`                            // IP
	Port      *int    `json:"Port,omitempty" binding:"omitempty,gt=0,lte=65535"` // 端口
	Username  *string `json:"Username,omitempty"`                                // 用户名
	Password  *string `json:"Password,omitempty"`                                // 密码
	GroupID   *int64  `json:"GroupID,omitempty"  binding:"omitempty,gt=0"`       // 组ID
	ProxyType *string `json:"ProxyType,omitempty"`                               // 代理类型，socks5
}

func (s *Svc) Update(params UpdateParams) error {
	// 先校验代理是否存在
	err := s.getRepo().GetByID(&models.Proxy{}, params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.NewErrorCode(common.ErrUpdateProxy, errors.New("代理不存在"))
		}
	}

	updateFields := map[string]interface{}{}
	if params.Port != nil {
		updateFields["port"] = *params.Port
	}
	if params.IP != nil {
		updateFields["ip"] = *params.IP
	}
	if params.Username != nil {
		updateFields["username"] = *params.Username
	}
	if params.Password != nil {
		updateFields["password"] = *params.Password
	}
	if params.GroupID != nil {
		groupAPI := group.NewGroupAPI(s.Ctx)
		hasGroup, err := groupAPI.CheckGroupID(*params.GroupID)
		if err != nil {
			return err
		}
		if !hasGroup {
			return errors.New("当前分组ID不是激活状态")
		}
		updateFields["group_id"] = *params.GroupID
	}
	if params.ProxyType != nil {
		updateFields["proxy_type"] = *params.ProxyType
	}

	err = s.getRepo().Update(params.ID, updateFields)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, "update proxy failed: %s", err.Error())
		return common.NewErrorCode(common.ErrUpdateProxy, err)
	}
	return nil
}

func (s *Svc) GetList(q models.GetListParams) (*common.ListData, error) {
	data := &common.ListData{}
	crud := s.getRepo()
	list := make([]models.Proxy, 0)

	total, err := crud.GetList(q, &models.Proxy{}, &list)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, "query list failed: %s", err.Error())
		return data, common.NewErrorCode(common.ErrGetList, err)
	}

	data.Counts = total
	data.Data = list
	return data, nil
}

func (s *Svc) Detail() (*models.Proxy, error) {
	p := &models.Proxy{}
	err := s.getRepo().GetByID(p, s.ID)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, "get proxy %d failed: %s", s.ID, err.Error())
		return nil, common.NewErrorCode(common.ErrGetDetail, err)
	}
	return p, nil
}

type DeleteParams struct {
	common.Test
	IPs []string `json:"IPs" binding:"required"`
}

func (s *Svc) Delete(params DeleteParams) error {
	err := s.getRepo().DeletesByIps(params.IPs)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, "delete proxies failed: %s", err.Error())
		return common.NewErrorCode(common.ErrDeleteGroup, err)
	}
	return nil
}

func (s *Svc) GetByIPs(ips []string) ([]models.Proxy, error) {
	if len(ips) == 0 {
		return nil, nil
	}

	query := models.GetListParams{
		IPs: ips,
	}

	var list []models.Proxy
	_, err := s.getRepo().GetList(query, &models.Proxy{}, &list)
	if err != nil {
		logger.ErrorfWithTrace(s.Ctx, "GetByIPs via GetList failed: %s", err.Error())
		return nil, err
	}

	return list, nil
}

func (s *Svc) GetByIP(ip string) (*models.Proxy, error) {
	proxy, err := s.getRepo().GetByIP(ip)
	if err != nil {
		return nil, err
	}
	return proxy, nil
}
