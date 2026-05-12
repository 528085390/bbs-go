// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"temp/post/api/internal/svc"
	"temp/post/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FeaturePostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFeaturePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FeaturePostLogic {
	return &FeaturePostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FeaturePostLogic) FeaturePost(req *types.ToggleReq) (resp *types.CommonResp, err error) {
	// todo: add your logic here and delete this line

	return
}
