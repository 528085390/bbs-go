package logic

import (
	"context"
	"errors"
	"fmt"
	"log"
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
	log.Println("userId: ", userId)
	log.Println("rolesSlice: ", rolesSlice)

	// 判断权限
	if !slices.Contains(rolesSlice, "admin") {
		return nil, errors.New("无权限")
	}

	// 删除板块
	res := l.svcCtx.Db.Where("id = ?", in.Id).Delete(&models.Section{})
	if res.RowsAffected == 0 {
		return nil, errors.New("板块不存在")
	}
	if res.Error != nil {
		return nil, res.Error
	}
	log.Printf(fmt.Sprintf("	管理员 %d 删除板块 %d 成功", userId, in.Id))

	return &rpc.DeleteSectionResponse{
		Success: true,
		Message: "删除板块成功",
	}, nil
}
