package logic

import (
	"context"
	"fmt"
	"mime"
	"path"
	"strings"
	"time"

	"temp/common/errs"
	"temp/common/errs/errorcode"
	"temp/common/httpctx"
	"temp/file/rpc/file"
	"temp/file/rpc/internal/svc"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetAvatarUploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAvatarUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAvatarUploadLogic {
	return &GetAvatarUploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAvatarUploadLogic) GetAvatarUploadURL(in *file.GetAvatarUploadRequest) (*file.UploadURLResponse, error) {
	userId, err := httpctx.GetUserId(l.ctx)
	if err != nil {
		logx.Errorf("get user id failed: %v", err)
		return nil, errs.Wrap(errorcode.Unauthorized, err, "获取用户信息失败")
	}

	fileExt := strings.ToLower(path.Ext(in.Filename))
	if !l.isValidFileType(fileExt) {
		logx.Errorf("invalid avatar file type: %s", fileExt)
		return nil, errs.New(errorcode.BadRequest, fmt.Sprintf("不支持的文件类型: %s", fileExt))
	}

	if in.MimeType != "" && !strings.Contains(in.MimeType, "/") {
		logx.Errorf("invalid mime type: %s", in.MimeType)
		return nil, errs.New(errorcode.BadRequest, fmt.Sprintf("无效的 mimeType 格式: %s", in.MimeType))
	}

	var contentType string
	if mimeType := strings.TrimSpace(in.MimeType); mimeType != "" {
		parsedType, _, err := mime.ParseMediaType(mimeType)
		if err != nil {
			logx.Errorf("invalid mime type: %s", in.MimeType)
			return nil, errs.New(errorcode.BadRequest, fmt.Sprintf("无效的 mimeType 格式: %s", in.MimeType))
		}
		if strings.HasPrefix(parsedType, "multipart/") {
			logx.Errorf("multipart upload not allowed: %s", parsedType)
			return nil, errs.New(errorcode.BadRequest, "不支持 multipart/form-data 上传，请使用 raw PUT 方式上传")
		}
		contentType = parsedType
	}

	// avatars/{userId}/{uuid}{ext}
	objectKey := fmt.Sprintf("avatars/%d/%s%s", userId, strings.ReplaceAll(uuid.New().String(), "-", ""), fileExt)

	timeout, err := time.ParseDuration(l.svcCtx.Config.OSS.UploadTimeout)
	if err != nil {
		logx.Errorf("parse upload timeout failed: %v", err)
		timeout = 15 * time.Minute
	}

	var signOptions []oss.Option
	if contentType != "" {
		signOptions = append(signOptions, oss.ContentType(contentType))
	}
	_, bucket, ok := l.svcCtx.OSS.Get()
	if !ok {
		return nil, errs.New(errorcode.ServerError, "OSS not ready")
	}
	signedURL, err := bucket.SignURL(objectKey, oss.HTTPPut, int64(timeout.Seconds()), signOptions...)
	if err != nil {
		logx.Errorf("生成签名 URL 失败: %v", err)
		return nil, errs.Wrap(errorcode.ServerError, err, "生成上传凭证失败")
	}

	downloadURL := l.buildDownloadURL(objectKey)

	logx.Infof("为用户 %d 生成头像上传凭证: %s", userId, objectKey)

	return &file.UploadURLResponse{
		UploadUrl:   signedURL,
		DownloadUrl: downloadURL,
		ObjectKey:   objectKey,
		ExpiresIn:   int32(timeout.Seconds()),
	}, nil
}

func (l *GetAvatarUploadLogic) isValidFileType(ext string) bool {
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true, ".bmp": true,
	}
	return allowedExts[ext]
}

func (l *GetAvatarUploadLogic) buildDownloadURL(objectKey string) string {
	return fmt.Sprintf("%s/%s", l.svcCtx.Config.OSS.BaseURL, objectKey)
}
