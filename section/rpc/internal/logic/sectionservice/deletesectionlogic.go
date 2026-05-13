package sectionservicelogic

import (
	"context"

	"temp/section/rpc/internal/svc"
	"temp/section/rpc/section/rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSectionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteSectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSectionLogic {
	return &DeleteSectionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteSectionLogic) DeleteSection(in *rpc.DeleteSectionRequest) (*rpc.DeleteSectionResponse, error) {
	// todo: add your logic here and delete this line

	return &rpc.DeleteSectionResponse{}, nil
}
