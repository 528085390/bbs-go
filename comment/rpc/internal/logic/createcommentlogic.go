package logic

import (
	"context"
	"fmt"
	"temp/common/errs"
	"temp/common/errs/errorcode"
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
		logx.Errorf("create comment invalid params: %v", err)
		return nil, err
	}

	// 帖子是否存在
	postRes, _ := l.svcCtx.PostService.ExistsPost(l.ctx, &post.IdPathReq{Id: postId})
	exists := postRes.Data
	if !exists {
		logx.Errorf("post not found: id=%d", postId)
		return nil, errs.New(errorcode.NotFound, fmt.Sprintf("帖子 %d 不存在", postId))
	}

	// 用户是否存在
	userRes, _ := l.svcCtx.UserService.ExistsUser(l.ctx, &user.IdRequest{Id: authorId})
	exists = userRes.Data
	if !exists {
		logx.Errorf("user not found: id=%d", authorId)
		return nil, errs.New(errorcode.ErrUserNotExist, fmt.Sprintf("用户 %d 不存在", authorId))
	}

	// 父评论是否存在
	var depth uint32
	var parentComment models.Comment
	if parentId != 0 {
		res := l.svcCtx.Db.Model(&models.Comment{}).Where("id = ?", parentId).First(&parentComment)
		if res.Error != nil {
			logx.Errorf("parent comment query failed: %v", res.Error)
			return nil, errs.Wrap(errorcode.NotFound, res.Error, "父评论不存在")
		}
		depth = parentComment.Depth + 1
	} else {
		depth = 0
	}

	// 校验评论深度
	if depth > 3 {
		logx.Errorf("comment depth exceeded: depth=%d", depth)
		return nil, errs.New(errorcode.BadRequest, "评论深度不能超过3")
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
		logx.Errorf("create comment failed: %v", res.Error)
		return nil, errs.Wrap(errorcode.ServerError, res.Error, "创建评论失败")
	}
	if res.RowsAffected == 0 {
		logx.Error("create comment affected rows is 0")
		return nil, errs.New(errorcode.ServerError, "创建评论失败")
	}

	logx.Infof("create comment success: id=%d post=%d author=%d", NewComment.ID, postId, authorId)

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
