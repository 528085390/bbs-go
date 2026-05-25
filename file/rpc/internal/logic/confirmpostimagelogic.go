package logic

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"temp/common/errs"
	"temp/common/errs/errorcode"
	"temp/common/httpctx"
	"temp/common/models"
	"temp/file/rpc/file"
	"temp/file/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfirmPostImageUploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfirmPostImageUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmPostImageUploadLogic {
	return &ConfirmPostImageUploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConfirmPostImageUploadLogic) ConfirmPostImageUpload(in *file.ConfirmPostImageUploadRequest) (*file.ConfirmUploadResponse, error) {
	userId, err := httpctx.GetUserId(l.ctx)
	if err != nil {
		logx.Errorf("get user id failed: %v", err)
		return nil, errs.Wrap(errorcode.Unauthorized, err, "获取用户信息失败")
	}

	_, bucket, ok := l.svcCtx.OSS.Get()
	if !ok {
		return nil, errs.New(errorcode.ServerError, "OSS not ready")
	}
	exists, err := bucket.IsObjectExist(in.ObjectKey)
	if err != nil {
		logx.Errorf("检查文件存在性失败: %v", err)
		return nil, errs.Wrap(errorcode.ServerError, err, "验证文件失败")
	}
	if !exists {
		logx.Errorf("object not exists: %s", in.ObjectKey)
		return nil, errs.New(errorcode.NotFound, "文件不存在，可能上传失败")
	}

	fileRecord := &models.File{
		ObjectKey: in.ObjectKey,
		Filename:  in.Filename,
		UserId:    userId,
		PostId:    in.PostId,
		Type:      "post_image",
		CreatedAt: time.Now(),
	}
	if err := l.svcCtx.DB.Model(&models.File{}).Create(fileRecord).Error; err != nil {
		logx.Errorf("保存文件记录失败: %v", err)
		return nil, errs.Wrap(errorcode.ServerError, err, "保存文件记录失败")
	}

	downloadURL := l.buildDownloadURL(in.ObjectKey)

	logx.Infof("用户 %d 上传帖子图片成功: post=%d %s (%d bytes)", userId, in.PostId, in.Filename, in.FileSize)

	return &file.ConfirmUploadResponse{
		FileId:   strconv.Itoa(int(fileRecord.ID)),
		Url:      downloadURL,
		Filename: in.Filename,
	}, nil
}

func (l *ConfirmPostImageUploadLogic) buildDownloadURL(objectKey string) string {
	return fmt.Sprintf("%s/%s", l.svcCtx.Config.OSS.BaseURL, objectKey)
}
