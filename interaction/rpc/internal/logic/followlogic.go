package logic

import (
	"context"
	"temp/common/errs"
	"temp/common/errs/errorcode"
	"temp/common/httpctx"
	"temp/common/models"
	"temp/common/valid"
	"temp/user/user"
	"time"

	"temp/interaction/rpc/interaction"
	"temp/interaction/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowLogic {
	return &FollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FollowLogic) Follow(in *interaction.FollowRequest) (*interaction.CommonResp, error) {
	// 参数校验
	followerId := in.FollowerId
	followingId := in.FollowingId
	err := valid.IsValidInt(followerId, followingId)
	if err != nil {
		logx.Errorf("follow invalid params: %v", err)
		return nil, err
	}
	followerRes, _ := l.svcCtx.UserRpc.ExistsUser(l.ctx, &user.IdRequest{Id: followerId})
	if !followerRes.Data {
		logx.Errorf("follow follower not found: id=%d", followerId)
		return nil, errs.New(errorcode.ErrUserNotExist, "用户不存在")
	}
	followingRes, _ := l.svcCtx.UserRpc.ExistsUser(l.ctx, &user.IdRequest{Id: followingId})
	if !followingRes.Data {
		logx.Errorf("follow following not found: id=%d", followingId)
		return nil, errs.New(errorcode.ErrUserNotExist, "用户不存在")
	}
	if followerId == followingId {
		logx.Error("follow self not allowed")
		return nil, errs.New(errorcode.BadRequest, "不能关注自己")
	}

	// 检验权限
	userId, err := httpctx.GetUserId(l.ctx)
	if err != nil {
		logx.Errorf("follow get user id failed: %v", err)
		return nil, err
	}
	if userId != followerId {
		logx.Errorf("follow forbidden: user=%d", followerId)
		return nil, errs.New(errorcode.Forbidden, "无权限")
	}

	// 处理关注
	op := in.Follow
	if op {
		follow := models.Follow{
			FollowerID:  followerId,
			FollowingID: followingId,
			CreatedAt:   time.Now(),
		}
		res := l.svcCtx.Db.Model(&models.Follow{}).Create(&follow)
		if res.Error != nil {
			logx.Errorf("follow create failed: %v", res.Error)
			return nil, errs.Wrap(errorcode.ServerError, res.Error, "关注失败")
		}
		logx.Infof("follow success: follower=%d following=%d", followerId, followingId)
		return &interaction.CommonResp{
			Message: "关注成功",
		}, nil
	} else {
		res := l.svcCtx.Db.Model(&models.Follow{}).Where("follower_id = ? and following_id = ?", followerId, followingId).Delete(&models.Follow{})
		if res.Error != nil {
			logx.Errorf("unfollow failed: %v", res.Error)
			return nil, errs.Wrap(errorcode.ServerError, res.Error, "取关失败")
		}
		logx.Infof("unfollow success: follower=%d following=%d", followerId, followingId)
		return &interaction.CommonResp{
			Message: "取关成功",
		}, nil
	}

}
