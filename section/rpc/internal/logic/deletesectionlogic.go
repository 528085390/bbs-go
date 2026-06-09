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

type DeleteSectionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteSectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSectionLogic {
	return &DeleteSectionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteSectionLogic) DeleteSection(in *rpc.DeleteSectionRequest) (*rpc.DeleteSectionResponse, error) {

	// 鉴权：仅 admin 可删除板块
	if err := httpctx.MustAdmin(l.ctx); err != nil {
		logx.Errorf("delete section forbidden: %v", err)
		return nil, err
	}

	// 删除板块
	res := l.svcCtx.Db.Where("id = ?", in.Id).Delete(&models.Section{})
	if res.RowsAffected == 0 {
		return nil, errs.New(errorcode.NotFound, "板块不存在")
	}
	if res.Error != nil {
		return nil, errs.Wrap(errorcode.ServerError, res.Error, "删除板块失败")
	}
	logx.Infof("删除板块 %d 成功", in.Id)

	return &rpc.DeleteSectionResponse{
		Success: true,
		Message: "删除板块成功",
	}, nil
}
