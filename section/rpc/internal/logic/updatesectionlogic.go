package logic

import (
	"context"
	"temp/section/rpc/internal/svc"
	"temp/section/rpc/section/rpc"

	"temp/common/models"

	"github.com/zeromicro/go-zero/core/logx"
	"k8s.io/kube-openapi/pkg/validation/errors"
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
		return nil, errors.New(400, "参数错误", nil)
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
		return nil, errors.New(500, "板块更新错误", res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, errors.New(500, "板块不存在", nil)
	}

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
