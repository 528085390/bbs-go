// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"temp/post/api/internal/svc"
	"temp/post/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PinPostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPinPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PinPostLogic {
	return &PinPostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PinPostLogic) PinPost(req *types.ToggleReq) (resp *types.CommonResp, err error) {
	// todo: add your logic here and delete this line

	return
}
