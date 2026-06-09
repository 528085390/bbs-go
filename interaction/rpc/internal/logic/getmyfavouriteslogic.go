package logic

import (
	"context"
	"temp/common/errs"
	"temp/common/errs/errorcode"
	"temp/common/httpctx"
	"temp/common/models"
	"temp/common/valid"

	"temp/interaction/rpc/interaction"
	"temp/interaction/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyFavouritesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMyFavouritesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyFavouritesLogic {
	return &GetMyFavouritesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMyFavouritesLogic) GetMyFavourites(in *interaction.GetMyFavouritesRequest) (*interaction.GetMyFavouritesResponse, error) {
	// 参数校验
	page := in.Page
	pageSize := in.PageSize
	userId, err := httpctx.GetUserId(l.ctx)
	if err != nil {
		logx.Errorf("get my favourites get user id failed: %v", err)
		return nil, errs.New(errorcode.Unauthorized, err)
	}
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	err = valid.IsValidInt(page, pageSize, userId)
	if err != nil {
		logx.Errorf("get my favourites invalid params: %v", err)
		return nil, err
	}

	// 获取收藏列表和总数
	var total int64
	var posts []models.Favorite
	res := l.svcCtx.Db.Model(&models.Favorite{}).Where("user_id = ?", userId).Count(&total)
	if res.Error != nil {
		logx.Errorf("get my favourites count failed: %v", res.Error)
		return nil, errs.Wrap(errorcode.ServerError, res.Error, "获取收藏数量失败")
	}
	res = l.svcCtx.Db.Model(&models.Favorite{}).Order("created_at DESC").Where("user_id = ?", userId).Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Find(&posts)
	if res.Error != nil {
		logx.Errorf("get my favourites list failed: %v", res.Error)
		return nil, errs.Wrap(errorcode.ServerError, res.Error, "获取收藏列表失败")
	}

	// 包装返回结果
	var postId []*interaction.PostId
	for _, post := range posts {
		postId = append(postId, &interaction.PostId{Id: post.PostID})
	}

	logx.Infof("get my favourites success: user=%d total=%d", userId, total)
	return &interaction.GetMyFavouritesResponse{
		Favourites: postId,
		Total:      total,
	}, nil
}
