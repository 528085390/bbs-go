package logic

import (
	"context"

	"temp/common/proto"
	"temp/search/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogic {
	return &SearchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Search 分页搜索
func (l *SearchLogic) Search(in *proto.SearchRequest) (*proto.SearchResponse, error) {
	res, err := l.svcCtx.PostRpc.Search(l.ctx, in)
	if err != nil {
		logx.Errorf("search rpc failed: %v", err)
		return nil, err
	}

	logx.Infof("search rpc success")
	return res, nil
}
