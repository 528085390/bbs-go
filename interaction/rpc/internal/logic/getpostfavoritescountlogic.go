package logic

import (
	"context"
	"temp/common/models"
	"temp/common/valid"
	"temp/post/rpc/post"

	"temp/interaction/rpc/interaction"
	"temp/interaction/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPostFavoritesCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPostFavoritesCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostFavoritesCountLogic {
	return &GetPostFavoritesCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPostFavoritesCountLogic) GetPostFavoritesCount(in *interaction.GetPostFavoritesCountRequest) (*interaction.GetPostFavoritesCountResponse, error) {
	postId := in.PostId
	// 参数校验
	err := valid.IsValidInt(postId)
	if err != nil {
		logx.Errorf("get post favorites count invalid params: %v", err)
		return nil, err
	}

	// 是否存在帖子
	_, err = l.svcCtx.PostRpc.ExistsPost(l.ctx, &post.IdPathReq{Id: postId})
	if err != nil {
		logx.Errorf("get post favorites count post not found: %v", err)
		return nil, err
	}

	// 获取收藏数
	var count int64
	res := l.svcCtx.Db.Model(&models.Favorite{}).Where("post_id = ?", postId).Count(&count)
	if res.Error != nil {
		logx.Errorf("get post favorites count failed: %v", res.Error)
		return nil, res.Error
	}

	logx.Infof("get post favorites count success: post=%d count=%d", postId, count)

	return &interaction.GetPostFavoritesCountResponse{
		Total: count,
	}, nil
}
