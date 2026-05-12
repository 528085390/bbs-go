package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
	"temp/common/auth"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. 定义不需要鉴权的路由（如登录、注册）
		publicPaths := []string{"/api/auth/login", "/api/auth/register"}
		for _, path := range publicPaths {
			if r.URL.Path == path {
				next.ServeHTTP(w, r)
				return
			}
		}

		// 2. 提取 Token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, `{"code": 401, "message": "unauthorized1"}`, http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// 3. 验证 Token
		ok, userId, roles, err := auth.ValidateAndParseToken("this is a secret", token)
		if !ok || err != nil {
			http.Error(w, `{"code": 401, "message": "unauthorized2"}`, http.StatusUnauthorized)
			return
		}
		// 4. 将用户信息存入 Context，供下游服务使用
		//ctx := r.Context()
		//ctx = context.WithValue(ctx, "userId", userId)
		//ctx = context.WithValue(ctx, "roles", roles)
		//next.ServeHTTP(w, r.WithContext(ctx))
		rolesBytes, _ := json.Marshal(roles)
		userIdBytes, _ := json.Marshal(userId)
		r.Header.Set("Grpc-Metadata-userid", string(userIdBytes))
		r.Header.Set("Grpc-Metadata-roles", string(rolesBytes))
		next.ServeHTTP(w, r)
	}
}
