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

type ListSectionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListSectionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSectionsLogic {
	return &ListSectionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListSectionsLogic) ListSections(in *rpc.Empty) (*rpc.SectionListResponse, error) {
	// 查询板块列表
	var sectionList []models.Section
	res := l.svcCtx.Db.Model(&models.Section{}).Find(&sectionList)
	if res.Error != nil {
		logx.Errorf("list sections failed: %v", res.Error)
		return nil, errs.Wrap(errorcode.ServerError, res.Error, "查询板块列表失败")
	}

	// 封装返回结果
	var sectionListResp []*rpc.SectionResponse
	for _, section := range sectionList {
		if !section.Visibility {
			continue
		}
		sectionListResp = append(sectionListResp, &rpc.SectionResponse{
			Id:          int64(section.ID),
			Title:       section.Title,
			Description: section.Description,
			OrderIndex:  int64(section.OrderIndex),
			Visibility:  section.Visibility,
			CreatedAt:   section.CreatedAt.String(),
			UpdatedAt:   section.UpdatedAt.String(),
		})
	}

	// 返回结果
	logx.Infof("list sections success: total=%d", len(sectionListResp))
	return &rpc.SectionListResponse{
		List:  sectionListResp,
		Total: int64(len(sectionListResp)),
	}, nil

}
