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

type GetPostImagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPostImagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostImagesLogic {
	return &GetPostImagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPostImagesLogic) GetPostImages(in *file.GetPostImagesRequest) (*file.GetPostImagesResponse, error) {
	var files []models.File
	res := l.svcCtx.DB.Where("post_id = ? AND type = ?", in.PostId, "post_image").Order("created_at asc").Find(&files)
	if res.Error != nil {
		logx.Errorf("get post images failed: %v", res.Error)
		return nil, errs.Wrap(errorcode.ServerError, res.Error, "查询失败")
	}

	resp := &file.GetPostImagesResponse{}
	for _, f := range files {
		url := fmt.Sprintf("%s/%s", l.svcCtx.Config.OSS.BaseURL, f.ObjectKey)
		resp.Items = append(resp.Items, &file.UrlResponse{Url: url})
	}

	logx.Infof("get post images success: post=%d count=%d", in.PostId, len(resp.Items))
	return resp, nil
}
