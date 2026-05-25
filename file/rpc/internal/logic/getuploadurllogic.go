package logic

import (
	"context"
	"fmt"
	"mime"
	"path"
	"strings"
	"temp/common/httpctx"
	"time"

	"temp/common/errs"
	"temp/common/errs/errorcode"
	"temp/file/rpc/file"
	"temp/file/rpc/internal/svc"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUploadURLLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUploadURLLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUploadURLLogic {
	return &GetUploadURLLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUploadURL 生成预签名上传 URL
func (l *GetUploadURLLogic) GetUploadURL(in *file.GetUploadURLRequest) (*file.UploadURLResponse, error) {
	// 1. 从 Context 中获取用户 ID（网关传递）
	userId, err := httpctx.GetUserId(l.ctx)
	if err != nil {
		logx.Errorf("get user id failed: %v", err)
		return nil, errs.Wrap(errorcode.Unauthorized, err, "获取用户信息失败")
	}

	// 2. 验证文件类型
	fileExt := strings.ToLower(path.Ext(in.Filename))
	if !l.isValidFileType(fileExt) {
		logx.Errorf("invalid file type: %s", fileExt)
		return nil, errs.New(errorcode.BadRequest, fmt.Sprintf("不支持的文件类型: %s", fileExt))
	}

	// 3. 验证 mimeType 格式，并禁止 multipart/form-data 这类表单上传
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

	// 4. 生成唯一的文件路径：uploads/{userId}/{uuid}{ext}
	objectKey := strings.ReplaceAll(uuid.New().String()+fileExt, "-", "")

	// 5. 解析上传超时时间
	timeout, err := time.ParseDuration(l.svcCtx.Config.OSS.UploadTimeout)
	if err != nil {
		logx.Errorf("parse upload timeout failed: %v", err)
		timeout = 15 * time.Minute // 默认 15 分钟
	}

	// 6. 生成预签名 URL（用于上传）
	// 注意：如果前端会发送 Content-Type，这里必须与实际请求保持一致；multipart/form-data 不适用
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

	// 7. 构建下载 URL
	downloadURL := l.buildDownloadURL(objectKey)

	logx.Infof("为用户 %d 生成上传凭证: %s", userId, objectKey)

	return &file.UploadURLResponse{
		UploadUrl:   signedURL,
		DownloadUrl: downloadURL,
		ObjectKey:   objectKey,
		ExpiresIn:   int32(timeout.Seconds()),
	}, nil
}

// isValidFileType 验证文件类型
func (l *GetUploadURLLogic) isValidFileType(ext string) bool {
	// 允许的图片类型
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true, ".bmp": true,
	}

	// 允许的文档类型
	allowedDocs := map[string]bool{
		".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true, ".ppt": true, ".pptx": true,
	}

	// 合并所有允许的类型
	allAllowed := make(map[string]bool)
	for k, v := range allowedExts {
		allAllowed[k] = v
	}
	for k, v := range allowedDocs {
		allAllowed[k] = v
	}

	return allAllowed[ext]
}

// buildDownloadURL 构建下载 URL
func (l *GetUploadURLLogic) buildDownloadURL(objectKey string) string {
	return fmt.Sprintf("%s/%s", l.svcCtx.Config.OSS.BaseURL, objectKey)
}
