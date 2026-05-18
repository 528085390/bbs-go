package logic

import (
	"context"
	"errors"
	"fmt"
	"temp/common/httpctx"
	"temp/common/models"
	"temp/common/valid"

	"temp/comment/rpc/comment"
	"temp/comment/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/exp/slices"
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
		return nil, err
	}

	// 查询评论
	var commentRes models.Comment
	res := l.svcCtx.Db.Model(&models.Comment{}).Where("id = ?", id).First(&commentRes)
	if res.Error != nil {
		return nil, res.Error
	}

	// 鉴权
	userId, err := httpctx.GetUserId(l.ctx)
	if err != nil {
		return nil, err
	}
	userRoles, err := httpctx.GetRoles(l.ctx)
	if err != nil {
		return nil, err
	}
	if commentRes.AuthorID != userId || !slices.Contains(userRoles, "admin") {
		return nil, errors.New(fmt.Sprintf("用户 %d 无权限删除评论 %d", userId, id))
	}

	// 删除评论
	res = l.svcCtx.Db.Model(&models.Comment{}).Where("id = ?", id).Update("content", "[该评论已删除]")
	if res.Error != nil {
		return nil, res.Error
	}

	// 返回
	return &comment.DeleteCommentResp{
		Success: true,
	}, nil
}
