package logic

import (
	"context"
	"temp/common/errs"
	"temp/common/errs/errorcode"
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
	sectionId := req.SectionId
	authorId, err := httpctx.GetUserId(l.ctx)
	if err != nil {
		logx.Errorf("get user id failed: %v", err)
		return nil, errs.New(errorcode.Unauthorized, err)
	}
	err = valid.IsValidString(title, content)
	if err != nil {
		logx.Errorf("create post invalid string params: %v", err)
		return nil, errs.New(errorcode.ParamError, err)
	}
	err = valid.IsValidInt(authorId, int64(sectionId))
	if err != nil {
		logx.Errorf("create post invalid int params: %v", err)
		return nil, errs.New(errorcode.ParamError, err)
	}
	user, err := l.svcCtx.UserRpc.GetUserInfoBy(l.ctx, &userclient.GetUserInfoRequest{
		Id:  authorId,
		Arg: "id",
	})
	if err != nil {
		logx.Errorf("get user info failed: %v", err)
		return nil, errs.New(errorcode.AuthorError, err)
	}
	if user == nil {
		logx.Errorf("user not found: id=%d", authorId)
		return nil, errs.New(errorcode.AuthorError, "用户不存在")
	}

	// 版块校验
	_, err = l.svcCtx.SectionRpc.GetSection(l.ctx, &sectionservice.GetSectionRequest{
		Id: sectionId,
	})
	if err != nil {
		logx.Errorf("get section failed: %v", err)
		return nil, errs.New(errorcode.SectionError, err)
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
		logx.Errorf("create post failed: %v", res.Error)
		return nil, errs.Wrap(res.Error)
	}

	logx.Infof("create post success: id=%d author=%d section=%d", newPost.ID, authorId, sectionId)

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
