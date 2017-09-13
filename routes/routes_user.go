package routes

import (
	"github.com/pagoda-tech/bastion/models"
	"github.com/pagoda-tech/macaron"
)

// UserShow 显示一个用户
func UserShow(ctx *macaron.Context, r APIRender, a Auth, db *models.DB) {
	// extract current user if 'current'
	id := uint(ctx.ParamsInt(":id"))
	if id == a.CurrentUser.ID {
		r.Success("user", a.CurrentUser)
		return
	}

	// check IsAdmin
	if !a.CanAccessUser(id) {
		r.Fail(UserNotFound, "没有找到该用户")
		return
	}

	// find
	u := &models.User{}
	db.First(u, id)

	if db.NewRecord(u) {
		r.Fail(UserNotFound, "没有找到该用户")
		return
	}

	r.Success("user", u)
}
