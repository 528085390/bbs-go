package logic

import (
	"context"
	"fmt"
	"strconv"
	"temp/common/httpctx"
	"temp/common/models"
	"temp/file/rpc/file"
	"temp/file/rpc/internal/svc"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"temp/common/errs"
	"temp/common/errs/errorcode"
)

type ConfirmUploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfirmUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmUploadLogic {
	return &ConfirmUploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ConfirmUpload 确认上传完成，保存文件记录到数据库
func (l *ConfirmUploadLogic) ConfirmUpload(in *file.ConfirmUploadRequest) (*file.ConfirmUploadResponse, error) {
	// 1. 从 Context 中获取用户 ID
	userId, err := httpctx.GetUserId(l.ctx)
	if err != nil {
		logx.Errorf("get user id failed: %v", err)
		return nil, errs.Wrap(errorcode.Unauthorized, err, "获取用户信息失败")
	}

	// 2. 验证文件是否存在于 OSS
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

	//3. TODO: 保存文件记录到数据库
	//这里需要集成 GORM 或其他 ORM
	fileRecord := &models.File{
		ObjectKey: in.ObjectKey,
		Filename:  in.Filename,
		UserId:    userId,
		CreatedAt: time.Now(),
	}
	if err := l.svcCtx.DB.Model(&models.File{}).Create(fileRecord).Error; err != nil {
		logx.Errorf("保存文件记录失败: %v", err)
		return nil, errs.Wrap(errorcode.ServerError, err, "保存文件记录失败")
	}

	// 4. 构建访问 URL
	downloadURL := l.buildDownloadURL(in.ObjectKey)

	logx.Infof("用户 %d 上传文件成功: %s (%d bytes)", userId, in.Filename, in.FileSize)

	return &file.ConfirmUploadResponse{
		FileId:   strconv.Itoa(int(fileRecord.ID)), // TODO: 替换为实际的数据库 ID
		Url:      downloadURL,
		Filename: in.Filename,
	}, nil
}

// buildDownloadURL 构建下载 URL
func (l *ConfirmUploadLogic) buildDownloadURL(objectKey string) string {
	return fmt.Sprintf("%s/%s", l.svcCtx.Config.OSS.BaseURL, objectKey)
}
