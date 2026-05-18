package logic

import (
	"context"
	"temp/common/models"
	"temp/user/internal/svc"
	"temp/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExistsUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewExistsUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExistsUserLogic {
	return &ExistsUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ExistsUserLogic) ExistsUser(in *user.IdRequest) (*user.ExistsResp, error) {

	userId := in.Id
	// 查询数据库
	res := l.svcCtx.Db.Model(&models.User{}).Where("id = ?", userId).First(&models.User{})
	if res.Error != nil {
		return &user.ExistsResp{
			Data: false,
		}, nil
	}

	// 用户存在
	return &user.ExistsResp{
		Data: true,
	}, nil

}
