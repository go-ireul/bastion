package routes

import (
	"fmt"
	"strings"

	"ireul.com/bastion/models"
	"ireul.com/web"
)

// Auth 认证结果
type Auth struct {
	Code         string
	Message      string
	CurrentToken *models.Token
	CurrentUser  *models.User
}

// SignedIn 是否已经登录
func (a Auth) SignedIn() bool {
	return a.CurrentToken != nil && a.CurrentUser != nil
}

// CanAccessUser 是否具有该用户的权限
func (a Auth) CanAccessUser(userID uint) bool {
	return a.SignedIn() && a.CurrentUser.ID == userID || a.CurrentUser.IsAdmin
}

// Authenticator 创建认证中间件
func Authenticator() interface{} {
	return func(ctx *web.Context, db *models.DB, r APIRender) {
		a := Auth{}

		k := extractBearer(ctx.Req)

		if len(k) > 0 {
			// find a Token
			t := &models.Token{}
			db.Where("secret = ?", k).Find(t)

			if t.ID > 0 {
				// find a User
				u := &models.User{}
				db.First(u, t.UserID)

				if u.ID > 0 && !u.IsBlocked {
					// touch token and user
					db.Touch(t)
					db.Touch(u)

					// assign user / token
					a.CurrentUser = u
					a.CurrentToken = t
				}
			}

			if a.CurrentUser == nil || a.CurrentToken == nil {
				a.Code = CredentialsInvalid
				a.Message = "无效的凭证"
			}
		} else {
			a.Code = CredentialsMissing
			a.Message = "没有凭证"
		}

		ctx.Map(a)
	}
}

// RequireAuth 检验认证结果
func RequireAuth() interface{} {
	return func(ctx *web.Context, a Auth, r APIRender) {
		if !a.SignedIn() {
			r.Fail(a.Code, a.Message)
		}
	}
}

// RequireServerToken require the server token should be set with Bearer Authorization
func RequireServerToken() interface{} {
	return func(ctx *web.Context, r APIRender, db *models.DB) {
		token := extractBearer(ctx.Req)

		if len(token) == 0 {
			r.Fail(CredentialsMissing, "credentials is missing")
			return
		}

		s := &models.Server{}
		db.Where("token = ?", token).First(s)

		if db.NewRecord(s) {
			r.Fail(CredentialsInvalid, "credentials is missing")
			return
		}

		ctx.Map(s)
	}
}

// RequireAdmin validates access user is admin
func RequireAdmin() interface{} {
	return func(ctx *web.Context, a Auth, r APIRender) {
		if !a.SignedIn() {
			r.Fail(a.Code, a.Message)
			return
		}
		if !a.CurrentUser.IsAdmin {
			r.Fail(PermissionDeny, "权限不足")
		}
	}
}

// ResolveCurrentUser 修正 current 为当前用户 ID
func ResolveCurrentUser(key string) interface{} {
	return func(ctx *web.Context, a Auth, r APIRender) {
		id := ctx.Params(key)
		if id == "current" {
			ctx.SetParams(key, fmt.Sprint(a.CurrentUser.ID))
		}
	}
}

// ResolveCurrentToken 修正 current 为当前 Token ID
func ResolveCurrentToken(key string) interface{} {
	return func(ctx *web.Context, a Auth, r APIRender) {
		id := ctx.Params(key)
		if id == "current" {
			ctx.SetParams(key, fmt.Sprint(a.CurrentToken.ID))
		}
	}
}

func extractBearer(req web.Request) (k string) {
	h := req.Header["Authorization"]
	if h != nil && len(h) > 0 {
		vs := strings.Split(strings.TrimSpace(h[len(h)-1]), " ")
		if len(vs) == 2 && vs[0] == "Bearer" {
			k = vs[1]
		}
	}
	return
}
