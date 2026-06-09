package logic

import (
	"context"
	"temp/common/errs"
	"temp/common/errs/errorcode"
	"temp/common/httpctx"
	"temp/common/models"
	"temp/common/valid"
	"temp/post/rpc/post"
	"time"

	"temp/interaction/rpc/interaction"
	"temp/interaction/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FavoriteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteLogic {
	return &FavoriteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FavoriteLogic) Favorite(in *interaction.FavoriteRequest) (*interaction.CommonResp, error) {
	// 参数校验
	postId := in.PostId
	op := in.Favorite
	userId, err := httpctx.GetUserId(l.ctx)
	if err != nil {
		logx.Errorf("favorite get user id failed: %v", err)
		return nil, errs.New(errorcode.Unauthorized, err)
	}
	err = valid.IsValidInt(userId, postId)
	if err != nil {
		logx.Errorf("favorite invalid params: %v", err)
		return nil, err
	}
	postRes, _ := l.svcCtx.PostRpc.ExistsPost(l.ctx, &post.IdPathReq{Id: postId})
	if !postRes.Data {
		logx.Errorf("favorite post not found: id=%d", postId)
		return nil, errs.New(errorcode.NotFound, "帖子不存在")
	}

	// 处理收藏
	if op {
		favorite := &models.Favorite{
			UserID:    userId,
			PostID:    postId,
			CreatedAt: time.Now(),
		}
		res := l.svcCtx.Db.Model(&models.Favorite{}).Create(&favorite)
		if res.Error != nil {
			logx.Errorf("favorite create failed: %v", res.Error)
			return nil, errs.Wrap(errorcode.ServerError, res.Error, "收藏失败")
		}
		logx.Infof("favorite success: user=%d post=%d", userId, postId)
		return &interaction.CommonResp{
			Message: "收藏成功",
		}, nil
	} else {
		res := l.svcCtx.Db.Model(&models.Favorite{}).Where("user_id = ? and post_id = ?", userId, postId).Delete(&models.Favorite{})
		if res.Error != nil {
			logx.Errorf("favorite delete failed: %v", res.Error)
			return nil, errs.Wrap(errorcode.ServerError, res.Error, "取消收藏失败")
		}
		logx.Infof("unfavorite success: user=%d post=%d", userId, postId)
		return &interaction.CommonResp{
			Message: "取消收藏成功",
		}, nil
	}

}
