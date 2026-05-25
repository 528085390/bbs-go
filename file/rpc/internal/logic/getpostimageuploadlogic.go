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
	"temp/file/rpc/file"
	"temp/file/rpc/internal/svc"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetPostImageUploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPostImageUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostImageUploadLogic {
	return &GetPostImageUploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPostImageUploadLogic) GetPostImageUploadURL(in *file.GetPostImageUploadRequest) (*file.UploadURLResponse, error) {
	// validate post_id
	if in.PostId == 0 {
		logx.Errorf("post_id is required")
		return nil, errs.New(errorcode.BadRequest, "post_id is required")
	}

	fileExt := strings.ToLower(path.Ext(in.Filename))
	// basic validation: allow image extensions only
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowedExts[fileExt] {
		logx.Errorf("invalid post image type: %s", fileExt)
		return nil, errs.New(errorcode.BadRequest, fmt.Sprintf("不支持的文件类型: %s", fileExt))
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

	objectKey := fmt.Sprintf("posts/%d/%s%s", in.PostId, strings.ReplaceAll(uuid.New().String(), "-", ""), fileExt)

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

	downloadURL := fmt.Sprintf("%s/%s", l.svcCtx.Config.OSS.BaseURL, objectKey)

	logx.Infof("生成帖子 %d 图片上传凭证: %s", in.PostId, objectKey)

	return &file.UploadURLResponse{
		UploadUrl:   signedURL,
		DownloadUrl: downloadURL,
		ObjectKey:   objectKey,
		ExpiresIn:   int32(timeout.Seconds()),
	}, nil
}
