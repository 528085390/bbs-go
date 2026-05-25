package logic

import (
	"context"
	"fmt"
	"strconv"
	"temp/common/errs"
	"temp/common/errs/errorcode"
	"temp/common/httpctx"
	"temp/common/models"

	"temp/post/rpc/internal/svc"
	"temp/post/rpc/post"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/exp/slices"
)

type FeaturePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFeaturePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FeaturePostLogic {
	return &FeaturePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FeaturePostLogic) FeaturePost(in *post.ToggleReq) (*post.CommonResp, error) {
	// 取参
	postId := in.Id
	val := in.Value
	userId, err := httpctx.GetUserId(l.ctx)
	if err != nil {
		logx.Errorf("get user id failed: %v", err)
		return nil, err
	}
	roles, err := httpctx.GetRoles(l.ctx)
	if err != nil {
		logx.Errorf("get roles failed: %v", err)
		return nil, err
	}

	// 权限校验
	if !slices.Contains(roles, "admin") {
		logx.Errorf("feature post forbidden: user=%d post=%d", userId, postId)
		return nil, errs.New(errorcode.Forbidden, fmt.Sprintf("用户 %d 无权限加精帖子 %d", userId, postId))
	}

	// 精华帖子
	var postRes models.Post
	res := l.svcCtx.Db.Model(&models.Post{}).Where("id = ?", postId).First(&postRes, postId)
	if res.Error != nil {
		logx.Errorf("get post failed: %v", res.Error)
		return nil, errs.Wrap(errorcode.ServerError, res.Error, "获取帖子失败")
	}
	res = l.svcCtx.Db.Model(&models.Post{}).Where("id = ?", postId).Update("featured", val)
	if res.Error != nil {
		logx.Errorf("update featured failed: %v", res.Error)
		return nil, errs.Wrap(errorcode.ServerError, res.Error, "更新加精失败")
	}

	logx.Infof("feature post success: id=%d value=%t", postId, val)

	return &post.CommonResp{
		Message: fmt.Sprintf("用户 %d 修改帖子 %d 精华状态为 %s 成功", userId, postId, strconv.FormatBool(val)),
	}, nil
}
