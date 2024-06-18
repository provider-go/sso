package main

import (
	"github.com/provider-go/pkg/logger"
	"github.com/provider-go/sso/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:123456@tcp(127.0.0.1:13306)/satellite?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	err = DB.AutoMigrate(models.SSOKey{})
	if err != nil {
		logger.Error(err)
	}

}
