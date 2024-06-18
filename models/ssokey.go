package models

import (
	"errors"
	"github.com/provider-go/sso/global"
	"gorm.io/gorm"
	"time"
)

type SSOKey struct {
	DID        string    `json:"did" gorm:"column:did;type:varchar(100);not null;default:'';primary_key;comment:'主键'"`
	PubKey     string    `json:"pubkey" gorm:"column:pubkey;type:varchar(100);not null;default:'';comment:用户公钥"`         // 用户公钥
	Status     int       `json:"status" gorm:"column:status;type:tinyint(1);not null;default:0;comment:账号状态:0(正常)1(禁用)"` // 账号状态:0(正常)1(禁用)
	CreateTime time.Time `json:"create_time" gorm:"autoCreateTime;comment:创建时间"`                                         // 创建时间
	UpdateTime time.Time `json:"update_time" gorm:"autoCreateTime;comment:更新时间"`                                         // 更新时间
}

func CreateSSOKey(did, pubkey string) error {
	return global.DB.Table("sso_keys").Create(&SSOKey{DID: did, PubKey: pubkey}).Error
}

func UpdateSSOKey(did, pubkey string) error {
	return global.DB.Table("sso_keys").Where("did = ?", did).Updates(map[string]interface{}{
		"pubkey": pubkey,
	}).Error
}

func DeleteSSOKey(did string) error {
	return global.DB.Table("sso_keys").Where("did = ?", did).Delete(&SSOKey{}).Error
}

func ListSSOKey(pageSize, pageNum int) ([]*SSOKey, int64, error) {
	var rows []*SSOKey
	//计算列表数量
	var count int64
	global.DB.Table("user_keys").Count(&count)

	if err := global.DB.Table("sso_keys").Order("create_time desc").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, count, nil
}

func ViewSSOKey(pubkey string) (*SSOKey, error) {
	row := new(SSOKey)
	if err := global.DB.Table("sso_keys").Where("pubkey = ?", pubkey).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 记录不存在的逻辑处理
			return nil, errors.New("ErrRecordNotFound")
		}
		return nil, err
	}
	return row, nil
}
