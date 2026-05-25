package logic

import (
	"context"
	"fmt"
	"temp/common/models"

	"temp/common/errs"
	"temp/common/errs/errorcode"
	"temp/file/rpc/file"
	"temp/file/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileLogic {
	return &GetFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFileLogic) GetFile(in *file.GetFileRequest) (*file.UrlResponse, error) {
	// TODO: 从数据库查询文件记录
	var fileRes models.File
	res := l.svcCtx.DB.Where("id = ?", in.Id).First(&fileRes)
	if res.Error != nil {
		logx.Errorf("get file failed: %v", res.Error)
		return nil, errs.New(errorcode.NotFound, "文件不存在")
	}

	logx.Infof("get file success: id=%s", in.Id)

	// 暂时返回示例 URL
	publicURL := fmt.Sprintf("https://%s.%s/%s", l.svcCtx.Config.OSS.BucketName, l.svcCtx.Config.OSS.Endpoint, fileRes.Filename)
	return &file.UrlResponse{
		Url: publicURL,
	}, nil
}
