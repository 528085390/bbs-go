// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"temp/common/models"
	"temp/section/api/internal/svc"
	"temp/section/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"k8s.io/kube-openapi/pkg/validation/errors"
)

type SectionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSectionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SectionListLogic {
	return &SectionListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SectionListLogic) SectionList() (resp *types.SectionListResp, err error) {
	// 查询板块列表
	var sectionList []models.Section
	res := l.svcCtx.Db.Find(&sectionList)
	if res.Error != nil {
		return nil, errors.New(500, "查询板块列表失败", nil)
	}

	// 封装返回结果
	var sectionListResp []types.Section
	for _, section := range sectionList {
		if !section.Visibility {
			continue
		}
		sectionListResp = append(sectionListResp, types.Section{
			Id:          int(section.ID),
			Title:       section.Title,
			Description: section.Description,
			OrderIndex:  section.OrderIndex,
			Visibility:  section.Visibility,
		})
	}
	return &types.SectionListResp{
		List:  sectionListResp,
		Total: int64(len(sectionListResp)),
	}, nil

}
