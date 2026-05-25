package logic

import (
	"context"
	"temp/common/models"
	"temp/common/valid"

	"temp/post/rpc/internal/svc"
	"temp/post/rpc/post"

	"github.com/zeromicro/go-zero/core/logx"
)

type MetaPostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMetaPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MetaPostLogic {
	return &MetaPostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MetaPostLogic) MetaPost(in *post.IdPathReq) (*post.PostMetaResp, error) {
	// 参数校验
	postId := in.Id
	err := valid.IsValidInt(postId)
	if err != nil {
		logx.Errorf("meta post invalid params: %v", err)
		return nil, err
	}

	// 查询
	var postRes models.Post
	res := l.svcCtx.Db.Model(&models.Post{}).Where("id = ?", postId).First(&postRes)
	if res.Error != nil {
		logx.Errorf("meta post query failed: %v", res.Error)
		return nil, res.Error
	}

	logx.Infof("meta post success: id=%d", postId)

	// 返回
	return &post.PostMetaResp{
		Id:        int64(postRes.ID),
		Title:     postRes.Title,
		AuthorId:  postRes.AuthorID,
		SectionId: uint64(postRes.SectionID),
	}, nil
}
