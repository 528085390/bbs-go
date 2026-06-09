package logic

import (
	"context"
	"temp/common/errs"
	"temp/common/errs/errorcode"
	"temp/common/httpctx"
	"temp/common/models"
	"temp/common/valid"
	"temp/post/rpc/post"

	"temp/interaction/rpc/interaction"
	"temp/interaction/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeLogic {
	return &LikeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LikeLogic) Like(in *interaction.LikeRequest) (*interaction.CommonResp, error) {
	// 参数校验
	postId := in.PostId
	op := in.Like
	userId, err := httpctx.GetUserId(l.ctx)
	if err != nil {
		logx.Errorf("like get user id failed: %v", err)
		return nil, errs.New(errorcode.Unauthorized, err)
	}
	err = valid.IsValidInt(userId, postId)
	if err != nil {
		logx.Errorf("like invalid params: %v", err)
		return nil, err
	}
	postRes, _ := l.svcCtx.PostRpc.ExistsPost(l.ctx, &post.IdPathReq{Id: postId})
	if !postRes.Data {
		logx.Errorf("like post not found: id=%d", postId)
		return nil, errs.New(errorcode.NotFound, "帖子不存在")
	}

	// 处理点赞
	if op {
		like := &models.Like{
			UserID: userId,
			PostID: postId,
		}
		res := l.svcCtx.Db.Model(&models.Like{}).Create(&like)
		if res.Error != nil {
			logx.Errorf("like create failed: %v", res.Error)
			return nil, errs.Wrap(errorcode.ServerError, res.Error, "点赞失败")
		}
	} else {
		res := l.svcCtx.Db.Model(&models.Like{}).Where("user_id = ? and post_id = ?", userId, postId).Delete(&models.Like{})
		if res.Error != nil {
			logx.Errorf("like delete failed: %v", res.Error)
			return nil, errs.Wrap(errorcode.ServerError, res.Error, "取消点赞失败")
		}
	}

	if op {
		logx.Infof("like success: user=%d post=%d", userId, postId)
	} else {
		logx.Infof("unlike success: user=%d post=%d", userId, postId)
	}

	return &interaction.CommonResp{
		Message: "ok",
	}, nil
}
