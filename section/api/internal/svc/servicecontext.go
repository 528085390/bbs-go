// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"temp/common/db"
	"temp/section/api/internal/config"

	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	Db     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Db:     db.GetDB(),
	}
}
