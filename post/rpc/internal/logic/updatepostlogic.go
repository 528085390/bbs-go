package logic

import (
	"context"
	"errors"
	"temp/common/httpctx"
	"temp/common/models"
	"temp/common/valid"
	"temp/section/rpc/sectionservice"
	"time"

	"temp/post/rpc/internal/svc"
	"temp/post/rpc/post"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePostLogic {
	return &UpdatePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdatePostLogic) UpdatePost(in *post.UpdatePostReq) (*post.PostResp, error) {
	// 参数校验
	postId := in.Id
	sectionId := in.SectionId
	title := in.Title
	content := in.Content
	authorId := in.AuthorId
	err := valid.IsValidString(title, content)
	err = valid.IsValidInt(postId, int64(sectionId), authorId)
	userId, err := httpctx.GetUserId(l.ctx)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.SectionRpc.GetSection(l.ctx, &sectionservice.GetSectionRequest{Id: sectionId})
	if err != nil {
		return nil, err
	}

	// 查询旧文章
	var PostRes models.Post
	res := l.svcCtx.Db.Model(&models.Post{}).Where("id = ?", postId).First(&PostRes)
	if res.Error != nil {
		return nil, res.Error
	}

	// 权限校验
	if PostRes.AuthorID != userId {
		return nil, errors.New("无权限修改该文章")
	}

	// 更新文章
	res = l.svcCtx.Db.Model(&models.Post{}).Where("id = ?", postId).Updates(models.Post{
		Title:     title,
		Content:   content,
		SectionID: int(sectionId),
	})
	if res.Error != nil {
		return nil, res.Error
	}

	// 返回结果
	return &post.PostResp{
		Id:        int64(PostRes.ID),
		Title:     title,
		Content:   content,
		AuthorId:  PostRes.AuthorID,
		SectionId: sectionId,
		Pinned:    PostRes.Pinned,
		Featured:  PostRes.Featured,
		CreatedAt: PostRes.CreatedAt.String(),
		UpdatedAt: time.Now().String(),
	}, nil
}
