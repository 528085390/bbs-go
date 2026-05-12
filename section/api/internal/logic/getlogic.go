// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"temp/section/api/internal/svc"
	"temp/section/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"k8s.io/kube-openapi/pkg/validation/errors"
)

type GetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLogic {
	return &GetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLogic) Get(req *types.GetSectionReq) (resp *types.GetSectionResp, err error) {
	var sectionId = req.Id
	var section types.Section
	res := l.svcCtx.Db.Table("sections").Where("id = ?", sectionId).First(&section)
	if res.Error != nil {
		return nil, errors.New(500, "查询板块信息失败", nil)
	}
	return &types.GetSectionResp{
		Item: section,
	}, nil

}
