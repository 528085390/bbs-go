package logic

import (
	"context"
	"errors"
	"temp/common/httpctx"
	"temp/common/models"
	"temp/common/valid"
	"temp/section/rpc/sectionservice"
	"temp/user/userclient"
	"time"

	"temp/post/rpc/internal/svc"
	"temp/post/rpc/post"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePostLogic {
	return &CreatePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreatePostLogic) CreatePost(req *post.PostRequest) (*post.PostResp, error) {
	// 参数校验
	title := req.Title
	content := req.Content
	authorId := req.AuthorId
	sectionId := req.SectionId
	userId, err := httpctx.GetUserId(l.ctx)
	if err != nil {
		return nil, err
	}
	err = valid.IsValidString(title, content)
	if err != nil {
		return nil, err
	}
	err = valid.IsValidInt(authorId, int64(sectionId))
	if err != nil {
		return nil, err
	}
	user, err := l.svcCtx.UserRpc.GetUserInfoBy(l.ctx, &userclient.GetUserInfoRequest{
		Id:  userId,
		Arg: "id",
	})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	// 权限校验
	if userId != authorId {
		return nil, errors.New("无权限发帖")
	}

	// 版块校验
	_, err = l.svcCtx.SectionRpc.GetSection(l.ctx, &sectionservice.GetSectionRequest{
		Id: sectionId,
	})
	if err != nil {
		return nil, err
	}

	// 创建帖子 添加帖子
	newPost := models.Post{
		Title:     title,
		Content:   content,
		AuthorID:  authorId,
		SectionID: int(sectionId),
		Pinned:    false,
		Featured:  false,
		ViewCount: 0,
	}
	res := l.svcCtx.Db.Model(&models.Post{}).Create(&newPost)
	if res.Error != nil {
		return nil, res.Error
	}

	// 封装返回数据
	return &post.PostResp{
		Id:        int64(newPost.ID),
		Title:     title,
		Content:   content,
		AuthorId:  authorId,
		SectionId: sectionId,
		Pinned:    false,
		Featured:  false,
		ViewCount: 0,
		CreatedAt: time.Now().String(),
		UpdatedAt: time.Now().String(),
	}, nil
}
