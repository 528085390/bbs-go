package sectionservicelogic

import (
	"context"

	"temp/section/rpc/internal/svc"
	"temp/section/rpc/section/rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSectionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateSectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSectionLogic {
	return &CreateSectionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateSectionLogic) CreateSection(in *rpc.CreateSectionRequest) (*rpc.CreateSectionResponse, error) {
	// todo: add your logic here and delete this line

	return &rpc.CreateSectionResponse{}, nil
}
