package errs

import (
	"context"
	"fmt"
	"temp/common/errs/errorcode"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const Key = "RpcError"

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}
		logx.Error(err)

		if rpcError, ok := From(err); ok {
			md := metadata.Pairs(Key, intToString(rpcError.Code))
			_ = grpc.SetTrailer(ctx, md)
			return nil, status.Error(codes.Internal, Message(rpcError))
		}

		_ = grpc.SetTrailer(ctx, metadata.Pairs(Key, intToString(errorcode.ServerError.Code)))
		return nil, status.Error(codes.Internal, err.Error())
	}
}

func intToString(v uint32) string {
	return fmt.Sprintf("%d", v)
}
