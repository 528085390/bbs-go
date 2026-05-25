package logic

import (
	"context"
	"temp/common/models"

	"temp/post/rpc/internal/svc"
	"temp/post/rpc/post"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExistsPostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewExistsPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExistsPostLogic {
	return &ExistsPostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ExistsPostLogic) ExistsPost(in *post.IdPathReq) (*post.ExistsResp, error) {
	postId := in.Id
	res := l.svcCtx.Db.Model(&models.Post{}).Where("id = ?", postId).First(&models.Post{})
	if res.Error != nil {
		logx.Errorf("exists post check failed: %v", res.Error)
		return &post.ExistsResp{
			Data: false,
		}, nil
	}

	logx.Infof("exists post success: id=%d", postId)
	return &post.ExistsResp{
		Data: true,
	}, nil
}
