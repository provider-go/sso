package global

import (
	"github.com/casbin/casbin/v2"
	"github.com/provider-go/pkg/cache"
	"gorm.io/gorm"
)

var (
	DB        *gorm.DB
	Cache     cache.Cache
	SecretKey string
	Casbin    *casbin.Enforcer
)
