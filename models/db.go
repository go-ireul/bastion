package models

import (
	"time"
	"github.com/pagoda-tech/bastion/utils"
	"github.com/pagoda-tech/gorm"
	"github.com/pagoda-tech/macaron"
)

// DB 封装 gorm.DB
type DB struct {
	*gorm.DB
}

// NewDB 创建一个新的 DB 实例
func NewDB(cfg *utils.Config) (db *DB, err error) {
	var db0 *gorm.DB
	if db0, err = gorm.Open("mysql", cfg.Database.URL); err != nil {
		return nil, err
	}
	// create
	db = &DB{db0}
	// enable log if dev
	if cfg.Bastion.Env == macaron.DEV {
		db.LogMode(true)
	}
	return
}

// AutoMigrate 自动执行数据库更新
func (db *DB) AutoMigrate() {
	db.DB.AutoMigrate(Token{}, Server{}, SSHKey{}, User{})
}

// Touch 更新一个模型的 UsedAt 字段
func (db *DB) Touch(m interface{}) {
	db.DB.Model(m).UpdateColumn("UsedAt", time.Now())
}
