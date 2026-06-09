package logic

import (
	"context"
	"fmt"
	"temp/common/errs"
	"temp/common/errs/errorcode"
	"temp/common/httpctx"
	"temp/common/models"
	"temp/common/valid"

	"temp/post/rpc/internal/svc"
	"temp/post/rpc/post"

	"github.com/zeromicro/go-zero/core/logx"
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
		logx.Errorf("delete post invalid params: %v", err)
		return nil, errs.New(errorcode.BadRequest, err)
	}
	userId, err := httpctx.GetUserId(l.ctx)
	if err != nil {
		logx.Errorf("get user id failed: %v", err)
		return nil, errs.New(errorcode.Unauthorized, err)
	}
	res := l.svcCtx.Db.Model(&models.Post{}).Where("id = ?", postId).First(&postRes)
	if res.Error != nil {
		logx.Errorf("query post failed: %v", res.Error)
		return nil, errs.Wrap(res.Error)
	}

	// 权限校验：作者本人或 admin 可删除
	if postRes.AuthorID != userId && !httpctx.IsAdmin(l.ctx) {
		logx.Errorf("delete post forbidden: user=%d post=%d", userId, postId)
		return nil, errs.New(errorcode.Forbidden, fmt.Sprintf("用户 %d 无权限删除帖子 %d", userId, postId))
	}

	// 删除帖子
	res = l.svcCtx.Db.Model(&models.Post{}).Where("id = ?", postId).Delete(&models.Post{})
	if res.Error != nil {
		logx.Errorf("delete post failed: %v", res.Error)
		return nil, errs.Wrap(res.Error)
	}

	logx.Infof("delete post success: id=%d user=%d", postId, userId)

	return &post.CommonResp{
		Message: fmt.Sprintf("用户 %d 删除帖子 %d 成功", userId, postId),
	}, nil
}
