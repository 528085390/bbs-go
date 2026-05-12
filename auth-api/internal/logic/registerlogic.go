// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"errors"
	"temp/common/auth"
	"temp/common/models"
	"temp/user/userclient"

	"temp/auth-api/internal/svc"
	"temp/auth-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	newUsername := req.Username
	newEmail := req.Email

	// 参数校验
	if newUsername == "" || newEmail == "" {
		resp := &types.RegisterResp{
			Message: "username or email is empty",
		}
		return resp, errors.New("username or email is empty")
	}

	// 重复校验
	var user *userclient.UserInfoResponse = nil
	user, err = l.svcCtx.UserRpc.GetUserInfoBy(l.ctx, &userclient.GetUserInfoRequest{
		Email: newEmail,
		Arg:   "email",
	})
	user, err = l.svcCtx.UserRpc.GetUserInfoBy(l.ctx, &userclient.GetUserInfoRequest{
		Username: newUsername,
		Arg:      "username",
	})
	if user != nil {
		logx.Error("username or email already exists")
		resp := &types.RegisterResp{
			Message: "username or email already exists",
		}
		return resp, errors.New("username or email already exists")
	}

	// 注册
	hashPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		logx.Error("HashPassword error")
		resp := &types.RegisterResp{
			Message: "register error",
		}
		return resp, errors.New("register error")
	}
	newUser := models.NewUser(
		req.Username,
		hashPassword,
		req.Email,
	)
	l.svcCtx.Db.Table("users").Create(&newUser)

	resp = &types.RegisterResp{
		Message: "register success",
	}
	return resp, nil

}
