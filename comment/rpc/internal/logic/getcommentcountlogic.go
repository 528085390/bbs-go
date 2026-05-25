package logic

import (
	"context"
	"temp/common/models"
	"temp/common/valid"

	"temp/comment/rpc/comment"
	"temp/comment/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentCountLogic {
	return &GetCommentCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetCommentCount 获取帖子评论数
func (l *GetCommentCountLogic) GetCommentCount(in *comment.GetCommentCountReq) (*comment.GetCommentCountResp, error) {
	// 参数校验
	postId := in.PostId
	err := valid.IsValidInt(postId)
	if err != nil {
		logx.Errorf("get comment count invalid params: %v", err)
		return nil, err
	}

	// 查询数据库
	var total int64
	res := l.svcCtx.Db.Model(&models.Comment{}).Where("post_id = ?", postId).Count(&total)
	if res.Error != nil {
		logx.Errorf("get comment count failed: %v", res.Error)
		return nil, err
	}

	logx.Infof("get comment count success: post=%d count=%d", postId, total)

	return &comment.GetCommentCountResp{
		Count: total,
	}, nil
}
