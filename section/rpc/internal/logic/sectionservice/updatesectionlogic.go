package sectionservicelogic

import (
	"context"

	"temp/section/rpc/internal/svc"
	"temp/section/rpc/section/rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSectionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateSectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSectionLogic {
	return &UpdateSectionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateSectionLogic) UpdateSection(in *rpc.UpdateSectionRequest) (*rpc.UpdateSectionResponse, error) {
	// todo: add your logic here and delete this line

	return &rpc.UpdateSectionResponse{}, nil
}
