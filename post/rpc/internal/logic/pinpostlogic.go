package logic

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"temp/common/httpctx"
	"temp/common/models"

	"temp/post/rpc/internal/svc"
	"temp/post/rpc/post"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/exp/slices"
)

type PinPostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPinPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PinPostLogic {
	return &PinPostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PinPostLogic) PinPost(in *post.ToggleReq) (*post.CommonResp, error) {
	// 取参
	postId := in.Id
	val := in.Value
	userId, err := httpctx.GetUserId(l.ctx)
	if err != nil {
		return nil, err
	}
	roles, err := httpctx.GetRoles(l.ctx)
	if err != nil {
		return nil, err
	}

	// 权限校验
	if !slices.Contains(roles, "admin") {
		return nil, errors.New(fmt.Sprintf("用户 %d 无权限帖子 %d", userId, postId))
	}

	// 置顶帖子
	var postRes models.Post
	res := l.svcCtx.Db.Model(&models.Post{}).Where("id = ?", postId).First(&postRes, postId)
	if res.Error != nil {
		return nil, res.Error
	}
	res = l.svcCtx.Db.Model(&models.Post{}).Where("id = ?", postId).Update("pinned", val)
	if res.Error != nil {
		return nil, res.Error
	}

	// 返回
	return &post.CommonResp{
		Message: fmt.Sprintf("用户 %d 改变帖子 %d 置顶状态为 %s 成功", userId, postId, strconv.FormatBool(val)),
	}, nil

}
