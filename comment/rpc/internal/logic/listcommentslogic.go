package logic

import (
	"context"
	"temp/comment/rpc/comment"
	"temp/comment/rpc/internal/svc"
	"temp/common/models"
	"temp/common/valid"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCommentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCommentsLogic {
	return &ListCommentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListComments 获取帖子评论列表
func (l *ListCommentsLogic) ListComments(in *comment.ListCommentsReq) (*comment.ListCommentsResp, error) {
	// 参数校验
	postId := in.PostId
	err := valid.IsValidInt(postId)
	if err != nil {
		return nil, err
	}

	// 查询评论列表
	var comments []*models.Comment
	res := l.svcCtx.Db.Model(&models.Comment{}).
		Where("post_id = ?", postId).
		Order("depth ASC, created_at DESC").
		Find(&comments)

	if res.Error != nil {
		return nil, res.Error
	}

	// 包装返回
	var commentsResp []*comment.CommentResp
	for _, c := range comments {
		commentsResp = append(commentsResp, &comment.CommentResp{
			Id:        int64(c.ID),
			PostId:    c.PostID,
			AuthorId:  c.AuthorID,
			ParentId:  c.ParentID,
			Content:   c.Content,
			Depth:     c.Depth,
			CreatedAt: c.CreatedAt.String(),
			UpdatedAt: c.UpdatedAt.String(),
		})
	}
	//sort.Slice(commentsResp, func(i, j int) bool {
	//	if commentsResp[i].Depth != commentsResp[j].Depth {
	//		return commentsResp[i].Depth < commentsResp[j].Depth
	//	} else {
	//		timei, _ := time.Parse("2006-01-02 15:04:05.999999 -0700 MST", commentsResp[i].CreatedAt)
	//		timej, _ := time.Parse("2006-01-02 15:04:05.999999 -0700 MST", commentsResp[j].CreatedAt)
	//		return timei.After(timej)
	//	}
	//
	//})
	return &comment.ListCommentsResp{
		Comments: commentsResp,
		Total:    uint64(len(commentsResp)),
	}, nil
}
