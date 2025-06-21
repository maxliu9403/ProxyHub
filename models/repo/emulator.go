package repo

import (
	"github.com/maxliu9403/ProxyHub/models"
	"github.com/maxliu9403/common/gormdb"
	"gorm.io/gorm"
	"time"
)

type EmulatorRepo interface {
	gormdb.GetByIDCrud
	GetList(q models.GetEmulatorListParams, model, list interface{}) (total int64, err error)
	Create(group *models.Emulator) error
	Update(uuid string, fields map[string]interface{}) error
	CreateBatch([]*models.Emulator) error
	DeletesByUuids(uuids []string) error
	GetByUuid(model interface{}, uuid string) error
	GetExistingUUIDs(uuids []string) ([]string, error)
	ListBriefByGroupID(groupID int64) ([]*models.EmulatorBrief, error)
	ListExpired(before time.Time) ([]*models.Emulator, error)
	DeletesByUuidsTx(tx *gorm.DB, uuids []string) error
}
