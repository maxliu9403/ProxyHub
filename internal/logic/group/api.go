package group

import (
	"context"

	"github.com/maxliu9403/common/gormdb"
)

type API interface {
	CheckGroupID(groupID int64) (bool, error)
}

type groupAPI struct {
	svc *Svc
}

func (g *groupAPI) CheckGroupID(groupID int64) (bool, error) {
	return g.svc.CheckGroupID(groupID)
}

func NewGroupAPI(ctx context.Context) API {
	return &groupAPI{
		svc: &Svc{
			Ctx: ctx,
			DB:  gormdb.Cli(ctx),
		},
	}
}
