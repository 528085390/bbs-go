package httpctx

import (
	"context"
	"encoding/json"
	"errors"

	"google.golang.org/grpc/metadata"
)

func getVal(ctx context.Context, key string) (res []byte, err error) {
	Vals := metadata.ValueFromIncomingContext(ctx, key)
	if len(Vals) == 0 {
		return nil, errors.New("请先登录 or 权限信息为空")
	}

	res = []byte(Vals[0])
	return
}

func GetUserId(ctx context.Context) (int64, error) {
	Vals, err := getVal(ctx, "gateway-userid")
	if err != nil {
		return 0, err
	}
	var userId float64
	err = json.Unmarshal(Vals, &userId)
	if err != nil {
		return 0, err
	}
	return int64(userId), nil
}

func GetRoles(ctx context.Context) ([]string, error) {
	Vals, err := getVal(ctx, "gateway-roles")
	if err != nil {
		return nil, err
	}
	var role []string
	err = json.Unmarshal(Vals, &role)
	if err != nil {
		return nil, err
	}
	return role, nil
}
