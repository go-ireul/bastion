package models

import (
	"github.com/jinzhu/gorm"
)

// Server 代表一个受管理的远端服务器
type Server struct {
	gorm.Model
}
