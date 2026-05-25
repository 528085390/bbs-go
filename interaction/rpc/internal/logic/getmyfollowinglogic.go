package logic

import (
	"context"
	"temp/common/errs"
	"temp/common/errs/errorcode"
	"temp/common/httpctx"
	"temp/common/models"
	"temp/common/valid"
	"temp/user/user"

	"temp/interaction/rpc/interaction"
	"temp/interaction/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyFollowingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMyFollowingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyFollowingLogic {
	return &GetMyFollowingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMyFollowingLogic) GetMyFollowing(in *interaction.GetMyFollowingRequest) (*interaction.GetMyFollowingResponse, error) {
	// 参数校验
	page := in.Page
	pageSize := in.PageSize
	userId := in.UserId
	err := valid.IsValidInt(page, pageSize, userId)
	if err != nil {
		logx.Errorf("get my following invalid params: %v", err)
		return nil, err
	}

	// 权限校验
	id, err := httpctx.GetUserId(l.ctx)
	if err != nil {
		logx.Errorf("get my following user id failed: %v", err)
		return nil, err
	}
	if id != userId {
		logx.Errorf("get my following forbidden: user=%d", userId)
		return nil, errs.New(errorcode.Forbidden, "无权限")
	}

	// 是否存在用户
	_, err = l.svcCtx.UserRpc.ExistsUser(l.ctx, &user.IdRequest{Id: userId})
	if err != nil {
		logx.Errorf("get my following user not found: %v", err)
		return nil, err
	}

	// 获取关注列表和总数
	var total int64
	var following []models.Follow
	res := l.svcCtx.Db.Model(&models.Follow{}).Where("follower_id = ?", userId).Count(&total)
	if res.Error != nil {
		logx.Errorf("get my following count failed: %v", res.Error)
		return nil, errs.Wrap(errorcode.ServerError, res.Error, "获取关注数量失败")
	}
	res = l.svcCtx.Db.Model(&models.Follow{}).Order("created_at DESC").Where("follower_id = ?", userId).Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Find(&following)
	if res.Error != nil {
		logx.Errorf("get my following list failed: %v", res.Error)
		return nil, errs.Wrap(errorcode.ServerError, res.Error, "获取关注列表失败")
	}

	// 包装返回结果
	var followingId []*interaction.UserId
	for _, f := range following {
		followingId = append(followingId, &interaction.UserId{Id: f.FollowingID})
	}

	logx.Infof("get my following success: user=%d total=%d", userId, total)
	return &interaction.GetMyFollowingResponse{
		Following: followingId,
		Total:     total,
	}, nil
}
