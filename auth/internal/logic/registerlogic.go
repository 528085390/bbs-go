package logic

import (
	"context"
	"temp/common/errs"
	"temp/common/errs/errorcode"
	"temp/common/models"
	"temp/common/password"
	"temp/user/userclient"

	"temp/auth/auth"
	"temp/auth/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *auth.RegisterReq) (*auth.RegisterResp, error) {
	newUsername := in.Username
	newEmail := in.Email

	// 参数校验
	if newUsername == "" || newEmail == "" {
		logx.Error("username or email is empty")
		return nil, errs.New(errorcode.BadRequest, "username or email is empty")
	}

	// 重复校验
	var user *userclient.UserInfoResponse = nil
	user, err := l.svcCtx.UserRpc.GetUserInfoBy(l.ctx, &userclient.GetUserInfoRequest{
		Email: newEmail,
		Arg:   "email",
	})
	if err != nil {
		logx.Errorf("check email exists failed: %v", err)
	}
	user, err = l.svcCtx.UserRpc.GetUserInfoBy(l.ctx, &userclient.GetUserInfoRequest{
		Username: newUsername,
		Arg:      "username",
	})
	if err != nil {
		logx.Errorf("check username exists failed: %v", err)
	}
	if user != nil {
		logx.Error("username or email already exists")
		return nil, errs.New(errorcode.ErrUserAlreadyExist, "username or email already exists")
	}

	// 注册
	hashPassword, err := password.HashPassword(in.Password)
	if err != nil {
		logx.Error("HashPassword error")
		return nil, errs.Wrap(errorcode.ServerError, err, "register error")
	}
	newUser := models.NewUser(in.Username, hashPassword, in.Email)
	if err := l.svcCtx.Db.Model(&models.User{}).Create(&newUser).Error; err != nil {
		logx.Errorf("create user failed: %v", err)
		return nil, errs.Wrap(errorcode.ServerError, err, "register error")
	}

	logx.Infof("register success: username=%s email=%s", newUsername, newEmail)
	return &auth.RegisterResp{Message: "register success"}, nil
}
