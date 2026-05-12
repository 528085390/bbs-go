package logic

import (
	"context"
	"temp/common/models"

	"temp/user/internal/svc"
	"temp/user/user"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetUserInfoByLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoByLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoByLogic {
	return &GetUserInfoByLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoByLogic) GetUserInfoBy(in *user.GetUserInfoRequest) (*user.UserInfoResponse, error) {
	arg := in.Arg

	// 参数校验 TODO err
	if arg == "" {
		logx.Error("GetUserInfoBy arg is empty")
		return nil, nil
	}
	if arg != "email" && arg != "username" && arg != "id" {
		logx.Error("GetUserInfoBy arg is not email or username or id")
		return nil, nil
	}

	var res *gorm.DB
	var resUser models.User
	if arg == "email" {
		res = l.svcCtx.Db.Table("users").Where("email = ?", in.Email).First(&resUser)
	} else if arg == "username" {
		res = l.svcCtx.Db.Table("users").Where("username = ?", in.Username).First(&resUser)
	} else {
		res = l.svcCtx.Db.Table("users").Where("id = ?", in.Id).First(&resUser)
	}
	if res != nil && res.Error != nil {
		logx.Errorf("GetUserInfoBy in database err: %v", res.Error)
		return nil, res.Error
	}

	return &user.UserInfoResponse{
		Id:          int64(resUser.ID),
		Username:    resUser.Username,
		Email:       resUser.Email,
		DisplayName: resUser.DisplayName,
		UpdatedAt:   resUser.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedAt:   resUser.CreatedAt.Format("2006-01-02 15:04:05"),
		Roles:       resUser.Roles,
	}, nil

}
