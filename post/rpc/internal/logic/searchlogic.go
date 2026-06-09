package logic

import (
	"context"
	"strings"
	"temp/comment/rpc/comment"
	"temp/common/errs"
	"temp/common/errs/errorcode"
	"temp/common/models"
	"temp/common/proto"
	"temp/interaction/rpc/interaction"
	"temp/post/rpc/internal/svc"
	"temp/user/userclient"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/exp/slices"
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

func (l *SearchLogic) Search(in *proto.SearchRequest) (*proto.SearchResponse, error) {
	key := strings.TrimSpace(in.Keyword)
	page := in.Page
	pageSize := in.PageSize
	sortField := strings.TrimSpace(in.SortField)
	sortOrder := strings.ToLower(strings.TrimSpace(in.SortOrder))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if sortField == "" {
		sortField = "created_at"
	}
	if sortOrder == "" {
		sortOrder = "desc"
	}

	// 参数校验
	ok := slices.Contains([]string{"created_at", "like_count", "view_count"}, sortField)
	if !ok {
		logx.Errorf("search invalid sort field: %s", sortField)
		return nil, errs.New(errorcode.BadRequest, "invalid sort field")
	}

	ok = slices.Contains([]string{"asc", "desc"}, sortOrder)
	if !ok {
		logx.Errorf("search invalid sort order: %s", sortOrder)
		return nil, errs.New(errorcode.BadRequest, "invalid sort order")
	}

	// 拼接排序
	orderBy := sortField + " " + strings.ToUpper(sortOrder)

	// 数据库查询
	query := l.svcCtx.Db.Model(&models.Post{})
	if key != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+key+"%", "%"+key+"%")
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		logx.Errorf("search count failed: %v", err)
		return nil, errs.Wrap(errorcode.ServerError, err, "search count failed")
	}

	// 获取帖子列表
	var posts []models.Post
	offset := int((page - 1) * pageSize)
	limit := int(pageSize)
	res := query.Order(orderBy).
		Offset(offset).
		Limit(limit).
		Find(&posts)
	if res.Error != nil {
		logx.Errorf("search list failed: %v", res.Error)
		return nil, errs.Wrap(errorcode.ServerError, res.Error, "search list failed")
	}

	// 封装搜索结果
	var searchResults []*proto.SearchResult

	// 获取作者，浏览量，评论数，点赞数等信息
	for _, p := range posts {

		r := proto.SearchResult{
			Id:             int64(p.ID),
			Title:          p.Title,
			AuthorId:       p.AuthorID,
			AuthorName:     "",
			ViewCount:      p.ViewCount,
			FavouriteCount: 0,
			CommentCount:   0,
			CreatedAt:      p.CreatedAt.String(),
			UpdatedAt:      p.UpdatedAt.String(),
		}

		// 获取作者信息
		userRes, err := l.svcCtx.UserRpc.GetUserInfoBy(l.ctx, &userclient.GetUserInfoRequest{
			Id:  p.AuthorID,
			Arg: "id",
		})
		if err != nil {
			logx.Errorf("get user info error: %v", err)
		} else if userRes != nil {
			r.AuthorName = userRes.Username
		}

		// 获取收藏信息
		interactRes, err := l.svcCtx.InteractionRpc.GetPostFavoritesCount(l.ctx, &interaction.GetPostFavoritesCountRequest{
			PostId: int64(p.ID),
		})
		if err != nil {
			logx.Errorf("get post favorites count error: %v", err)
		} else if interactRes != nil {
			r.FavouriteCount = interactRes.Total
		}

		// 获取评论数
		commentRes, err := l.svcCtx.CommentRpc.GetCommentCount(l.ctx, &comment.GetCommentCountReq{
			PostId: int64(p.ID),
		})
		if err != nil {
			logx.Errorf("get comments count error: %v", err)
		} else if commentRes != nil {
			r.CommentCount = commentRes.Count
		}

		searchResults = append(searchResults, &r)

	}

	// 返回结果
	logx.Infof("search success: key=%s total=%d", key, total)
	return &proto.SearchResponse{
		Total: total,
		Items: searchResults,
	}, nil

}
