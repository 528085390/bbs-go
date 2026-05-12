package logic

import (
	"context"
	"log"
	"temp/common/models"
	"temp/section/rpc/internal/svc"
	"temp/section/rpc/section/rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"k8s.io/kube-openapi/pkg/validation/errors"
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
	log.Printf("title: %s, description: %s, orderIndex: %d, visibility: %t", title, description, orderIndex, visibility)
	if title == "" || description == "" || orderIndex < 0 {
		return nil, errors.New(400, "参数错误", nil)
	}

	// 判断板块是否存在
	res := l.svcCtx.Db.Table("sections").Where("title = ?", title).First(&models.Section{})
	if res.Error == nil {
		return nil, errors.New(400, "板块已存在", res.Error)
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
		return nil, errors.New(500, "板块创建错误", res.Error)
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
	return &rpc.CreateSectionResponse{
		Section: &resp,
	}, nil

}
