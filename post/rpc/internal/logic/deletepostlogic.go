package logic

import (
	"context"
	"errors"
	"fmt"
	"temp/common/httpctx"
	"temp/common/models"
	"temp/common/valid"

	"temp/post/rpc/internal/svc"
	"temp/post/rpc/post"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/exp/slices"
)

type DeletePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePostLogic {
	return &DeletePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeletePostLogic) DeletePost(in *post.IdPathReq) (*post.CommonResp, error) {
	// 参数校验
	var postRes models.Post
	postId := in.Id
	err := valid.IsValidInt(postId)
	if err != nil {
		return nil, err
	}
	userId, err := httpctx.GetUserId(l.ctx)
	if err != nil {
		return nil, err
	}
	roles, err := httpctx.GetRoles(l.ctx)
	if err != nil {
		return nil, err
	}
	res := l.svcCtx.Db.Model(&models.Post{}).Where("id = ?", postId).First(&postRes)
	if res.Error != nil {
		return nil, res.Error
	}

	// 权限校验
	if postRes.AuthorID != userId && !slices.Contains(roles, "admin") {
		return nil, errors.New(fmt.Sprintf("用户 %d 无权限删除帖子 %d", userId, postId))
	}

	// 删除帖子
	res = l.svcCtx.Db.Model(&models.Post{}).Where("id = ?", postId).Delete(&models.Post{})
	if res.Error != nil {
		return nil, res.Error
	}

	return &post.CommonResp{
		Message: fmt.Sprintf("用户 %d 删除帖子 %d 成功", userId, postId),
	}, nil
}
