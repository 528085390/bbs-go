package logic

import (
	"context"
	"fmt"
	"strconv"
	"temp/common/errs"
	"temp/common/errs/errorcode"
	"temp/common/httpctx"
	"temp/common/models"

	"temp/post/rpc/internal/svc"
	"temp/post/rpc/post"

	"github.com/zeromicro/go-zero/core/logx"
)

type PinPostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPinPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PinPostLogic {
	return &PinPostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PinPostLogic) PinPost(in *post.ToggleReq) (*post.CommonResp, error) {
	// 取参
	postId := in.Id
	val := in.Value
	// 鉴权：仅 admin 可置顶
	if err := httpctx.MustAdmin(l.ctx); err != nil {
		logx.Errorf("pin post forbidden: post=%d", postId)
		return nil, err
	}

	// 置顶帖子
	var postRes models.Post
	res := l.svcCtx.Db.Model(&models.Post{}).First(&postRes, postId)
	if res.Error != nil {
		logx.Errorf("get post failed: %v", res.Error)
		return nil, errs.Wrap(errorcode.ServerError, res.Error, "获取帖子失败")
	}
	res = l.svcCtx.Db.Model(&models.Post{}).Where("id = ?", postId).Update("pinned", val)
	if res.Error != nil {
		logx.Errorf("update pinned failed: %v", res.Error)
		return nil, errs.Wrap(errorcode.ServerError, res.Error, "更新置顶失败")
	}

	logx.Infof("pin post success: id=%d value=%t", postId, val)

	// 返回
	return &post.CommonResp{
		Message: fmt.Sprintf("帖子 %d 置顶状态改为 %s 成功", postId, strconv.FormatBool(val)),
	}, nil

}
