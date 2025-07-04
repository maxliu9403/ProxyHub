/*
@Date: 2022/4/15 16:33
@Author: max.liu
@File : repo
@Desc:
*/

package factory

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm/clause"

	"github.com/maxliu9403/ProxyHub/models"
	"github.com/maxliu9403/ProxyHub/models/repo"
	"github.com/maxliu9403/common/gadget"
	"github.com/maxliu9403/common/gormdb"
	"github.com/maxliu9403/common/rsql"
	"gorm.io/gorm"
)

type proxyCrudImpl struct {
	Conn *gorm.DB
}

func ProxyRepo(db *gorm.DB) repo.ProxyRepo {
	return &proxyCrudImpl{Conn: db}
}

func (r *proxyCrudImpl) GetList(q models.GetListParams, model, list interface{}) (total int64, err error) {
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
	if len(q.IPs) > 0 {
		db.Where("ip IN ?", q.IPs)
	}

	if len(q.Ports) > 0 {
		db.Where("port IN ?", q.Ports)
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

func (r *proxyCrudImpl) GetByID(model interface{}, id int64) error {
	crud := gormdb.NewCRUD(r.Conn)
	err := crud.GetByID(model, id)
	return err
}

func (r *proxyCrudImpl) Deletes(ids []int64) (err error) {
	err = r.Conn.Delete(&models.Proxy{}, ids).Error
	return err
}

func (r *proxyCrudImpl) Create(proxy *models.Proxy) error {
	return r.Conn.Create(proxy).Error
}

func (r *proxyCrudImpl) Update(id int64, fields map[string]interface{}) error {
	return r.Conn.Model(&models.Proxy{}).Where("id = ?", id).Updates(fields).Error
}

func (r *proxyCrudImpl) CreateBatch(proxies []*models.Proxy) error {
	if len(proxies) == 0 {
		return nil
	}

	ips := make([]string, 0, len(proxies))
	for _, e := range proxies {
		ips = append(ips, e.IP)
	}

	// 只物理删除已软删除（delete_time 不为空）的冲突 UUID 记录，mysql唯一索引会有冲突
	if err := r.Conn.
		Unscoped().
		Where("ip IN ?", ips).
		Where("delete_time IS NOT NULL").
		Delete(&models.Proxy{}).Error; err != nil {
		return err
	}

	// 插入新记录
	return r.Conn.Create(&proxies).Error
}

func (r *proxyCrudImpl) DeletesByIps(IPs []string) error {
	return r.Conn.Model(&models.Proxy{}).
		Where("ip IN ?", IPs).
		Update("delete_time", time.Now()).Error
}

func (r *proxyCrudImpl) GetByIP(ip string) (*models.Proxy, error) {
	proxy := &models.Proxy{}
	err := r.Conn.Where("ip = ?", ip).First(proxy).Error
	return proxy, err
}

func (r *proxyCrudImpl) IncrementInUseTx(tx *gorm.DB, ip string, count int) error {
	return tx.Model(&models.Proxy{}).
		Where("ip = ?", ip).
		UpdateColumn("inuse_count", gorm.Expr("inuse_count + ?", count)).Error
}

// DecrementInUseTx 在事务中安全递减某个 IP 的 inuse_count
func (r *proxyCrudImpl) DecrementInUseTx(tx *gorm.DB, ip string, count int) error {
	// 加锁读取
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("ip = ?", ip).First(&models.Proxy{}).Error; err != nil {
		return err
	}

	// 原子更新
	return tx.Model(&models.Proxy{}).
		Where("ip = ?", ip).
		UpdateColumn("inuse_count", gorm.Expr("GREATEST(inuse_count - ?, 0)", count)).Error
}

// GetByIPForUpdate 查询指定 IP 并加锁，事务中使用
func (r *proxyCrudImpl) GetByIPForUpdate(ip string) (*models.Proxy, error) {
	var proxy models.Proxy
	err := r.Conn.Raw(`SELECT * FROM tbl_proxy WHERE ip = ? FOR UPDATE`, ip).Scan(&proxy).Error
	if err != nil {
		return nil, err
	}
	return &proxy, nil
}

func (r *proxyCrudImpl) ListByGroupID(groupID int64) ([]*models.ProxyBrief, error) {
	var proxies []*models.ProxyBrief
	err := r.Conn.Model(&models.Proxy{}).
		Select("ip,port,username,password").
		Where("group_id = ?", groupID).
		Scan(&proxies).Error
	return proxies, err
}
