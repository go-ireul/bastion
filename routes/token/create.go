package token

import (
	"github.com/pagoda-tech/bastion/models"
	"github.com/pagoda-tech/bastion/routes/middlewares"
	"github.com/pagoda-tech/bastion/utils"
	"github.com/pagoda-tech/macaron"
	"time"
)

// CreateForm 创建 Token 表单
type CreateForm struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Desc     string `json:"desc"`
}

// Create 创建 Token 路由
func Create(ctx *macaron.Context, db *models.DB, f CreateForm, r *middlewares.Render) {
	u := &models.User{}
	// find
	db.Where("login = ?", f.Login).First(u)
	if u.ID == 0 {
		r.Fail("login_failed", "用户名或密码错误")
		return
	}
	// check password
	if !u.CheckPassword(f.Password) {
		r.Fail("login_failed", "用户名或密码错误")
		return
	}
	// create token
	t := &models.Token{
		UserID: u.ID,
		Desc:   f.Desc,
	}
	t.GenerateSecret()
	db.Create(t)
	// touch user
	db.Model(u).Update("UsedAt", time.Now())
	// return data
	r.Success(func(m utils.Map) {
		m.Set("token", utils.NewMap(
			"id",
			t.ID,
			"secret",
			t.Secret,
		))
	})
}
