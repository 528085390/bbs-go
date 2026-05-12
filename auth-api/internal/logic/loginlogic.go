// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"errors"
	"temp/common/auth"
	"temp/common/models"

	"temp/auth-api/internal/svc"
	"temp/auth-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// 参数校验
	username := req.Username
	password := req.Password
	if username == "" || password == "" {
		logx.Error("username or password is empty")
		return nil, errors.New("username or password is empty")
	}

	// 查询用户
	var user models.User

	res := l.svcCtx.Db.Table("users").Where("username = ?", username).First(&user)
	if res.Error != nil {
		logx.Errorf("find user in database err: %v", res.Error)
		return nil, res.Error
	}

	if !auth.CheckPasswordHash(password, user.Password) {
		logx.Error("password is incorrect")
		return nil, errors.New("password is incorrect")
	}

	// 生成token
	token, err := auth.GenerateAccessToken(l.svcCtx.Config.JwtSecret, int64(user.ID), user.Roles)
	if err != nil {
		logx.Error("generate access token error")
		return nil, err
	}
	return &types.LoginResp{
		Token: token,
	}, nil
}
