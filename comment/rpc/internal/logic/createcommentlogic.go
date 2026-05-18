package logic

import (
	"context"
	"errors"
	"fmt"
	"temp/common/models"
	"temp/common/valid"
	"temp/post/rpc/post"
	"temp/user/user"
	"time"

	"temp/comment/rpc/comment"
	"temp/comment/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCommentLogic {
	return &CreateCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateComment 创建评论
func (l *CreateCommentLogic) CreateComment(in *comment.CreateCommentReq) (*comment.CommentResp, error) {
	// 参数校验
	postId := in.PostId
	authorId := in.AuthorId
	parentId := in.ParentId
	err := valid.IsValidInt(postId, authorId)
	if err != nil {
		return nil, err
	}

	// 帖子是否存在
	postRes, _ := l.svcCtx.PostService.ExistsPost(l.ctx, &post.IdPathReq{Id: postId})
	exists := postRes.Data
	if !exists {
		return nil, errors.New(fmt.Sprintf("帖子 %d 不存在", postId))
	}

	// 用户是否存在
	userRes, _ := l.svcCtx.UserService.ExistsUser(l.ctx, &user.IdRequest{Id: authorId})
	exists = userRes.Data
	if !exists {
		return nil, errors.New(fmt.Sprintf("用户 %d 不存在", authorId))
	}

	// 父评论是否存在
	var depth uint32
	var parentComment models.Comment
	if parentId != 0 {
		res := l.svcCtx.Db.Model(&models.Comment{}).Where("id = ?", parentId).First(&parentComment)
		if res.Error != nil {
			return nil, res.Error
		}
		depth = parentComment.Depth + 1
	} else {
		depth = 0
	}

	// 校验评论深度
	if depth > 3 {
		return nil, errors.New("评论深度不能超过3")
	}

	// 创建评论
	NewComment := models.Comment{
		PostID:   postId,
		AuthorID: authorId,
		ParentID: parentId,
		Content:  in.Content,
		Depth:    depth,
	}

	// 插入数据库
	res := l.svcCtx.Db.Model(&models.Comment{}).Create(&NewComment)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, errors.New("创建评论失败")
	}

	// 返回结果
	return &comment.CommentResp{
		Id:        int64(NewComment.ID),
		PostId:    NewComment.PostID,
		AuthorId:  NewComment.AuthorID,
		ParentId:  NewComment.ParentID,
		Content:   NewComment.Content,
		Depth:     NewComment.Depth,
		CreatedAt: time.Now().String(),
		UpdatedAt: time.Now().String(),
	}, nil

}
