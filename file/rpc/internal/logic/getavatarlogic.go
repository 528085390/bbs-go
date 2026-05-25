package logic

import (
	"context"
	"fmt"

	"temp/common/errs"
	"temp/common/errs/errorcode"
	"temp/common/models"
	"temp/file/rpc/file"
	"temp/file/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAvatarLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAvatarLogic {
	return &GetAvatarLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAvatarLogic) GetAvatar(in *file.GetAvatarRequest) (*file.UrlResponse, error) {
	var f models.File
	res := l.svcCtx.DB.Where("user_id = ? AND type = ?", in.UserId, "avatar").Order("created_at desc").First(&f)
	if res.Error != nil {
		logx.Errorf("get avatar failed: %v", res.Error)
		return nil, errs.New(errorcode.NotFound, "头像不存在")
	}

	logx.Infof("get avatar success: user=%d", in.UserId)

	publicURL := fmt.Sprintf("%s/%s", l.svcCtx.Config.OSS.BaseURL, f.ObjectKey)
	return &file.UrlResponse{Url: publicURL}, nil
}
