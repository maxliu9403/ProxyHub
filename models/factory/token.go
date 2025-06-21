package factory

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/maxliu9403/ProxyHub/models"
	"github.com/maxliu9403/ProxyHub/models/repo"
	"github.com/maxliu9403/common/gadget"
	"github.com/maxliu9403/common/gormdb"
	"github.com/maxliu9403/common/rsql"
	"gorm.io/gorm"
)

type tokenCrudImpl struct {
	Conn *gorm.DB
}

func TokenRepo(db *gorm.DB) repo.TokenRepo {
	return &tokenCrudImpl{Conn: db}
}

func (r *tokenCrudImpl) GetList(q models.GetTokenListParams, model, list interface{}) (total int64, err error) {
	db := r.Conn.Model(model)

	// 指定字段
	if len(q.Fields) > 0 {
		db.Select(strings.Join(q.Fields, ", "))
	}

	parseColumnFunc := func(s string) string { return r.Conn.NamingStrategy.ColumnName("", s) }

	// 精确字段模糊匹配
	if len(q.FuzzyField) > 0 {
		for k, v := range q.FuzzyField {
			columnName := parseColumnFunc(k)
			db.Scopes(gormdb.KeywordGenerator([]string{columnName}, v))
		}
	}

	// 全局模糊
	if q.Keyword != "" {
		fields := gadget.FieldsFromModel(model, db, true).GetStringField()
		db.Scopes(gormdb.KeywordGenerator(fields, q.Keyword))
	}

	if len(q.GroupIDs) > 0 {
		db.Where("group_id IN ?", q.GroupIDs)
	}

	// 自定义查询条件
	if q.Query != "" {
		// 把传递过来的Query字段通过gorm的字段命名策略转义成数据库字段
		preParser, e := rsql.NewPreParser(rsql.MysqlPre(parseColumnFunc))
		if e != nil {
			err = e
			return total, err
		}

		preStmt, values, err := preParser.ProcessPre(q.Query)
		if err != nil {
			return total, err
		}

		db.Where(preStmt, values...)
	}

	// 排序
	if q.Order != "" {
		orderList := strings.Split(q.Order, ",")

		for _, o := range orderList {
			orderKey := strings.Split(o, " ")
			switch len(orderKey) {
			case 1:
				columnName := parseColumnFunc(orderKey[0])
				db.Order(columnName)
			case 2:
				columnName := parseColumnFunc(orderKey[0])
				order := strings.ToUpper(orderKey[1])
				if order != "DESC" && order != "ASC" {
					order = "ASC"
				}

				db.Order(fmt.Sprintf("%s %s", columnName, order))
			}
		}
	}

	// 计数
	db = db.Count(&total)

	// 分页
	if q.Limit > 0 && q.Offset >= 0 {
		db.Limit(q.Limit).Offset(q.Offset)
	}

	err = db.Find(list).Error

	return total, err
}

func (r *tokenCrudImpl) Create(token *models.Token) error {
	return r.Conn.Create(token).Error
}

func (r *tokenCrudImpl) Deletes(token []string) error {
	return r.Conn.Delete(&models.Token{}, "token in ?", token).Error
}

func (r *tokenCrudImpl) Get(token string) (*models.Token, error) {
	var t models.Token
	err := r.Conn.First(&t, "token = ?", token).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *tokenCrudImpl) IsValid(token string) (bool, error) {
	// 查询 token 是否存在
	t, err := r.Get(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Token 不存在，返回 false
			return false, errors.New("token不存在")
		}
		// 查询失败（数据库错误等）
		return false, err
	}

	// 永久有效
	if t.ExpireAt == nil {
		return true, nil
	}

	// 检查是否已过期
	if time.Now().Before(*t.ExpireAt) {
		return true, nil
	}

	// 已过期
	return false, nil
}

func (r *tokenCrudImpl) GetValidTokensByGroup(groupID int64, now time.Time) ([]models.Token, error) {
	var list []models.Token
	query := r.Conn.Model(&models.Token{}).Where("group_id = ?", groupID)
	query = query.Where("expire_at IS NULL OR expire_at > ?", now)

	err := query.Find(&list).Error
	return list, err
}

// GetByGroupID 查询指定 GroupID 的 Token（主键为 token，非 group_id）
func (r *tokenCrudImpl) GetByGroupID(groupID int64) (*models.Token, error) {
	var token models.Token
	err := r.Conn.Model(&token).
		Where("group_id = ?", groupID).
		First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}
