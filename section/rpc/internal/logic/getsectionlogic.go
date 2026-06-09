package logic

import (
	"context"
	"temp/common/errs"
	"temp/common/errs/errorcode"
	"temp/common/models"

	"temp/section/rpc/internal/svc"
	"temp/section/rpc/section/rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSectionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSectionLogic {
	return &GetSectionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSectionLogic) GetSection(req *rpc.GetSectionRequest) (*rpc.SectionResponse, error) {
	var sectionId = req.Id
	var section models.Section
	res := l.svcCtx.Db.Model(&models.Section{}).Where("id = ?", sectionId).First(&section)
	if res.Error != nil {
		logx.Errorf("get section failed: %v", res.Error)
		return nil, errs.Wrap(errorcode.NotFound, res.Error, "板块不存在")
	}

	logx.Infof("get section success: id=%d", section.ID)
	return &rpc.SectionResponse{
		Id:          int64(section.ID),
		Title:       section.Title,
		Description: section.Description,
		OrderIndex:  int64(section.OrderIndex),
		Visibility:  section.Visibility,
		CreatedAt:   section.CreatedAt.String(),
		UpdatedAt:   section.UpdatedAt.String(),
	}, nil
}
