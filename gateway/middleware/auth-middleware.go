package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
	"temp/common/errs/errorcode"
	"temp/common/response"
	"temp/common/tokenUtil"
	"temp/gateway/config"
)

// matchPath 支持路径参数匹配（如 /api/posts/123 匹配 /api/posts/:id）
func matchPath(actualPath, patternPath string) bool {
	actualParts := strings.Split(strings.Trim(actualPath, "/"), "/")
	patternParts := strings.Split(strings.Trim(patternPath, "/"), "/")

	if len(actualParts) != len(patternParts) {
		return false
	}

	for i, patternPart := range patternParts {
		// 如果模式部分是 :param 形式，则匹配任何值
		if strings.HasPrefix(patternPart, ":") {
			continue
		}
		// 否则必须完全匹配
		if actualParts[i] != patternPart {
			return false
		}
	}

	return true
}

func AuthMiddleware(c config.Config) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 检查是否为公开路由（支持路径参数匹配）
			isPublic := false
			for _, route := range c.PublicRoutes {
				if strings.EqualFold(r.Method, route.Method) && matchPath(r.URL.Path, route.Path) {
					isPublic = true
					break
				}
			}

			if isPublic {
				next.ServeHTTP(w, r)
				return
			}

			// 2. 提取 Token
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				writeError(w, errorcode.Unauthorized)
				return
			}
			token := strings.TrimPrefix(authHeader, "Bearer ")

			// 3. 验证 Token
			ok, userId, roles, err := tokenUtil.ValidateAndParseToken(c.JwtSecret, token)
			if !ok || err != nil {
				writeError(w, errorcode.Unauthorized)
				return
			}
			// 4. 将用户信息存入 Context，供下游服务使用
			//ctx := r.Context()
			//ctx = context.WithValue(ctx, "userId", userId)
			//ctx = context.WithValue(ctx, "roles", roles)
			//next.ServeHTTP(w, r.WithContext(ctx))
			rolesBytes, err := json.Marshal(roles)
			if err != nil {
				writeError(w, errorcode.ServerError)
				return
			}
			userIdBytes, err := json.Marshal(userId)
			if err != nil {
				writeError(w, errorcode.ServerError)
				return
			}
			r.Header.Set("Grpc-Metadata-userid", string(userIdBytes))
			r.Header.Set("Grpc-Metadata-roles", string(rolesBytes))
			next.ServeHTTP(w, r)
		}
	}
}

func writeError(w http.ResponseWriter, errCode errorcode.ErrorCode) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// 实际使用 json.Marshal
	json.NewEncoder(w).Encode(response.Error(int(errCode.Code), errCode.Msg, nil))
}
