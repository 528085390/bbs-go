package logic

import (
	"context"

	"temp/common/errs"
	"temp/common/errs/errorcode"
	"temp/common/httpctx"
	"temp/common/models"
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
	id := in.Id
	title := in.Title
	description := in.Description
	orderIndex := int(in.OrderIndex)
	visibility := in.Visibility

	// 参数校验
	if title == "" || description == "" || orderIndex < 0 || id < 0 {
		logx.Error("update section invalid params")
		return nil, errs.New(errorcode.BadRequest, "参数错误")
	}

	// 鉴权：仅 admin 可更新板块
	if err := httpctx.MustAdmin(l.ctx); err != nil {
		logx.Errorf("update section forbidden: %v", err)
		return nil, err
	}

	//封装新板块
	newSection := models.Section{
		Title:       title,
		Description: description,
		OrderIndex:  orderIndex,
		Visibility:  visibility,
	}

	res := l.svcCtx.Db.Model(&models.Section{}).Where("id = ?", id).Select("title", "description", "order_index", "visibility").Updates(&newSection)
	if res.Error != nil {
		logx.Errorf("update section failed: %v", res.Error)
		return nil, errs.Wrap(errorcode.ServerError, res.Error, "板块更新错误")
	}
	if res.RowsAffected == 0 {
		logx.Errorf("update section not found: id=%d", id)
		return nil, errs.New(errorcode.NotFound, "板块不存在")
	}

	logx.Infof("update section success: id=%d", id)
	return &rpc.UpdateSectionResponse{
		Section: &rpc.SectionResponse{
			Id:          in.Id,
			Title:       newSection.Title,
			Description: newSection.Description,
			OrderIndex:  int64(newSection.OrderIndex),
			Visibility:  newSection.Visibility,
		},
	}, nil
}
