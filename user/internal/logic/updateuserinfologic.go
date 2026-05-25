package logic

import (
	"context"
	"fmt"
	"temp/common/errs"
	"temp/common/errs/errorcode"
	"temp/common/models"
	"time"

	"temp/user/internal/svc"
	"temp/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(req *user.GetUserInfoRequest) (*user.UserInfoResponse, error) {
	// find user && get old info
	userRes := models.User{}
	res := l.svcCtx.Db.Table("users").Find(&userRes, req.Id)
	if res.Error != nil {
		logx.Error("find user in database error")
		return nil, errs.Wrap(res.Error)
	}
	if res.RowsAffected == 0 {
		logx.Error("user not found")
		return nil, errs.New(errorcode.UserNotFound, fmt.Sprintf("userId: %d", req.Id))
	}

	// update user
	if req.Username != "" {
		userRes.Username = req.Username
	}
	if req.Email != "" {
		userRes.Email = req.Email
	}
	if req.DisplayName != "" {
		userRes.DisplayName = req.DisplayName
	}
	userRes.UpdatedAt = time.Now()

	res = l.svcCtx.Db.Table("users").Save(&userRes)
	if res.Error != nil {
		logx.Errorf("update user failed: %v", res.Error)
		return nil, errs.Wrap(res.Error)
	}

	logx.Infof("update user success: id=%d", userRes.ID)

	return &user.UserInfoResponse{
		Id:          int64(userRes.ID),
		Username:    userRes.Username,
		Email:       userRes.Email,
		DisplayName: userRes.DisplayName,
		UpdatedAt:   userRes.UpdatedAt.String(),
		CreatedAt:   userRes.CreatedAt.String(),
		Roles:       userRes.Roles,
	}, nil

}
