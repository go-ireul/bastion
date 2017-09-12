package models

import (
	"github.com/pagoda-tech/bastion/utils"
	"github.com/pagoda-tech/gorm"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"
)

// UserLoginRegexp 用户登录名正则表达式
var UserLoginRegexp = regexp.MustCompile(`^[a-zA-Z]\w{3,15}$`)

// UserPasswordMinLen 用户密码最短长度
const UserPasswordMinLen = 6

// User 代表一个用户
type User struct {
	gorm.Model

	// Login 登录名
	Login string `gorm:"unique_index" json:"login"`
	// Nickname 昵称
	Nickname string `json:"nickname"`
	// PasswordDigest bcrypt 加密后的密码
	PasswordDigest string `gorm:"type:text" json:"-"`
	// IsBlocked 用户是否被禁用
	IsBlocked bool `gorm:"not null" json:"isBlocked"`
	// IsAdmin 用户是否是管理员
	IsAdmin bool `gorm:"not null" json:"isAdmin"`

	// Fingerprint Sandbox 公钥指纹
	Fingerprint string `gorm:"unique_index" json:"fingerprint"`
	// PublicKey Sandbox 公钥
	PublicKey string `gorm:"type:text" json:"-"`
	// PrivateKey Sandbox 私钥
	PrivateKey string `gorm:"type:text" json:"-"`
	// UsedAt 最后一次使用时间
	UsedAt *time.Time `json:"usedAt"`
}

// GenerateSSHKey 为该用户填充一个新的 SSH 密钥对
func (u *User) GenerateSSHKey() (err error) {
	ou := *u
	if u.Fingerprint, u.PublicKey, u.PrivateKey, err = utils.GenerateSSHKeyPair(); err != nil {
		// failed, restore u
		*u = ou
	}
	return
}

// SetPassword 为该用户更新密码
func (u *User) SetPassword(p string) (err error) {
	var b []byte
	if b, err = bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost); err != nil {
		return
	}
	u.PasswordDigest = string(b)
	return
}

// CheckPassword 检查密码
func (u *User) CheckPassword(p string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordDigest), []byte(p)) == nil
}
