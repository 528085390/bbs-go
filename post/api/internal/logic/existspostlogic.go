// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"temp/post/api/internal/svc"
	"temp/post/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExistsPostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExistsPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExistsPostLogic {
	return &ExistsPostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExistsPostLogic) ExistsPost(req *types.IdPathReq) (resp *types.ExistsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
