package logic

import (
	"context"
	"temp/common/errs"
	"temp/common/errs/errorcode"
	"temp/common/models"
	"temp/common/tokenUtil"

	"temp/auth/auth"
	"temp/auth/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *auth.LoginReq) (*auth.LoginResp, error) {
	// 参数校验
	username := in.Username
	password := in.Password
	if username == "" || password == "" {
		logx.Error("username or password is empty")
		return nil, errs.New(errorcode.BadRequest, "username or password is empty")
	}

	// 查询用户
	var user models.User
	res := l.svcCtx.Db.Table("users").Where("username = ?", username).First(&user)
	if res.Error != nil {
		logx.Errorf("find user in database err: %v", res.Error)
		return nil, errs.Wrap(errorcode.ServerError, res.Error, "find user failed")
	}

	if !tokenUtil.CheckPasswordHash(password, user.Password) {
		logx.Error("password is incorrect")
		return nil, errs.New(errorcode.ErrPasswordIncorrect, "password is incorrect")
	}

	// 生成token
	token, err := tokenUtil.GenerateAccessToken(l.svcCtx.Config.JwtSecret, int64(user.ID), user.Roles)
	if err != nil {
		logx.Error("generate access token error")
		return nil, errs.Wrap(errorcode.ServerError, err, "generate access token error")
	}

	logx.Infof("login success: userId=%d username=%s", user.ID, username)
	return &auth.LoginResp{
		Token: token,
	}, nil
}
