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

	"github.com/maxliu9403/ProxyHub/models"
	"github.com/maxliu9403/ProxyHub/models/repo"
	"github.com/maxliu9403/common/gadget"
	"github.com/maxliu9403/common/gormdb"
	"github.com/maxliu9403/common/rsql"
	"gorm.io/gorm"
)

type emulatorCrudImpl struct {
	Conn *gorm.DB
}

func EmulatorRepo(db *gorm.DB) repo.EmulatorRepo {
	return &emulatorCrudImpl{Conn: db}
}

func (r *emulatorCrudImpl) GetList(q models.GetEmulatorListParams, model, list interface{}) (total int64, err error) {
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

	if q.GroupIDs != nil {
		db.Where("group_id IN ?", q.GroupIDs)
	}

	if q.UUIDS != nil {
		db.Where("uuid IN ?", q.UUIDS)
	}

	if q.BrowserIDs != nil {
		db.Where("browser_id IN ?", q.BrowserIDs)
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

func (r *emulatorCrudImpl) GetByID(model interface{}, id int64) error {
	crud := gormdb.NewCRUD(r.Conn)
	err := crud.GetByID(model, id)
	return err
}

func (r *emulatorCrudImpl) Create(emulator *models.Emulator) error {
	return r.Conn.Create(emulator).Error
}

func (r *emulatorCrudImpl) Update(uuid string, fields map[string]interface{}) error {
	return r.Conn.Model(&models.Emulator{}).Where("uuid = ?", uuid).Updates(fields).Error
}

func (r *emulatorCrudImpl) CreateBatch(emulator []*models.Emulator) error {
	return r.Conn.Create(&emulator).Error
}

func (r *emulatorCrudImpl) DeletesByUuids(uuids []string) error {
	return r.Conn.Model(&models.Emulator{}).
		Where("uuid IN ?", uuids).
		Update("delete_time", time.Now()).Error
}

func (r *emulatorCrudImpl) GetByUuid(model interface{}, uuid string) error {
	return r.Conn.Model(&models.Emulator{}).Where("uuid = ?", uuid).First(&model).Error
}

func (r *emulatorCrudImpl) GetExistingUUIDs(uuids []string) ([]string, error) {
	var exist []string
	err := r.Conn.
		Model(&models.Emulator{}).
		Where("uuid IN ?", uuids).
		Pluck("uuid", &exist).Error
	return exist, err
}

func (r *emulatorCrudImpl) ListBriefByGroupID(groupID int64) ([]*models.EmulatorBrief, error) {
	var list []*models.EmulatorBrief
	err := r.Conn.Model(&models.Emulator{}).
		Select("browser_id, uuid, ip").
		Where("group_id = ?", groupID).
		Scan(&list).Error
	return list, err
}

func (r *emulatorCrudImpl) ListExpired(before time.Time) ([]*models.Emulator, error) {
	var list []*models.Emulator
	err := r.Conn.Model(&models.Emulator{}).
		Where("update_time < ? AND delete_time IS NULL AND ip != ''", before).
		Find(&list).Error
	return list, err
}

func (r *emulatorCrudImpl) DeletesByUuidsTx(tx *gorm.DB, uuids []string) error {
	return tx.Model(&models.Emulator{}).
		Where("uuid IN ?", uuids).
		Update("delete_time", time.Now()).Error
}
