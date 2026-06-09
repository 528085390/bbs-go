package logic

import (
	"context"
	"fmt"
	"temp/common/errs"
	"temp/common/errs/errorcode"
	"temp/common/httpctx"
	"temp/common/models"
	"temp/common/valid"

	"temp/comment/rpc/comment"
	"temp/comment/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCommentLogic {
	return &DeleteCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteComment 删除评论
func (l *DeleteCommentLogic) DeleteComment(in *comment.DeleteCommentReq) (*comment.DeleteCommentResp, error) {
	// 参数校验
	id := in.Id
	err := valid.IsValidInt(id)
	if err != nil {
		logx.Errorf("delete comment invalid params: %v", err)
		return nil, err
	}

	// 查询评论
	var commentRes models.Comment
	res := l.svcCtx.Db.Model(&models.Comment{}).Where("id = ?", id).First(&commentRes)
	if res.Error != nil {
		logx.Errorf("query comment failed: %v", res.Error)
		return nil, errs.Wrap(errorcode.NotFound, res.Error, "评论不存在")
	}

	// 鉴权：作者本人或 admin 可删除
	userId, err := httpctx.GetUserId(l.ctx)
	if err != nil {
		logx.Errorf("get user id failed: %v", err)
		return nil, errs.Wrap(errorcode.Unauthorized, err, "获取用户信息失败")
	}
	if commentRes.AuthorID != userId && !httpctx.IsAdmin(l.ctx) {
		logx.Errorf("delete comment forbidden: user=%d comment=%d", userId, id)
		return nil, errs.New(errorcode.Forbidden, fmt.Sprintf("用户 %d 无权限删除评论 %d", userId, id))
	}

	// 删除评论
	res = l.svcCtx.Db.Model(&models.Comment{}).Where("id = ?", id).Update("content", "[该评论已删除]")
	if res.Error != nil {
		logx.Errorf("delete comment update failed: %v", res.Error)
		return nil, errs.Wrap(errorcode.ServerError, res.Error, "删除评论失败")
	}

	logx.Infof("delete comment success: id=%d user=%d", id, userId)

	// 返回
	return &comment.DeleteCommentResp{
		Success: true,
	}, nil
}
