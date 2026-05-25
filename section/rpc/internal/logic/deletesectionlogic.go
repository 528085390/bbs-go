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
	"golang.org/x/exp/slices"
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

	// 获取权限信息
	userId, err := httpctx.GetUserId(l.ctx)
	rolesSlice, err := httpctx.GetRoles(l.ctx)
	if err != nil {
		return nil, err
	}
	logx.Infof("userId: %d", userId)
	logx.Infof("rolesSlice: %v", rolesSlice)

	// 判断权限
	if !slices.Contains(rolesSlice, "admin") {
		return nil, errs.New(errorcode.Forbidden, "无权限")
	}

	// 删除板块
	res := l.svcCtx.Db.Where("id = ?", in.Id).Delete(&models.Section{})
	if res.RowsAffected == 0 {
		return nil, errs.New(errorcode.NotFound, "板块不存在")
	}
	if res.Error != nil {
		return nil, errs.Wrap(errorcode.ServerError, res.Error, "删除板块失败")
	}
	logx.Infof("管理员 %d 删除板块 %d 成功", userId, in.Id)

	return &rpc.DeleteSectionResponse{
		Success: true,
		Message: "删除板块成功",
	}, nil
}
