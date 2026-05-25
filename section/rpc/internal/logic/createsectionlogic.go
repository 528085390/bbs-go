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
	// 参数校验
	title := in.Title
	description := in.Description
	orderIndex := int(in.OrderIndex)
	visibility := in.Visibility
	logx.Infof("title: %s, description: %s, orderIndex: %d, visibility: %t", title, description, orderIndex, visibility)
	if title == "" || description == "" || orderIndex < 0 {
		return nil, errs.New(errorcode.BadRequest, "参数错误")
	}

	// 判断板块是否存在
	res := l.svcCtx.Db.Table("sections").Where("title = ?", title).First(&models.Section{})
	if res.Error == nil {
		logx.Errorf("section already exists: title=%s", title)
		return nil, errs.New(errorcode.BadRequest, "板块已存在")
	}

	// 创建板块
	var newSection = models.Section{
		Title:       title,
		Description: description,
		OrderIndex:  orderIndex,
		Visibility:  visibility,
	}
	res = l.svcCtx.Db.Table("sections").Create(&newSection)
	if res.Error != nil {
		logx.Errorf("create section failed: %v", res.Error)
		return nil, errs.Wrap(errorcode.ServerError, res.Error, "创建板块失败")
	}

	// 返回结果
	var resp = rpc.SectionResponse{
		Id:          int64(newSection.ID),
		Title:       newSection.Title,
		Description: newSection.Description,
		OrderIndex:  int64(newSection.OrderIndex),
		Visibility:  newSection.Visibility,
		CreatedAt:   newSection.CreatedAt.String(),
		UpdatedAt:   newSection.UpdatedAt.String(),
	}
	logx.Infof("create section success: id=%d title=%s", newSection.ID, newSection.Title)
	return &rpc.CreateSectionResponse{
		Section: &resp,
	}, nil

}
