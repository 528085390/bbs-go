package errorcode

type ErrorCode struct {
	Code uint32
	Msg  string
}

var (
	ServerError  = ErrorCode{Code: 10000, Msg: "服务器错误"}
	Unauthorized = ErrorCode{Code: 10001, Msg: "未授权"}
	Forbidden    = ErrorCode{Code: 10002, Msg: "无权限"}
	NotFound     = ErrorCode{Code: 10003, Msg: "资源不存在"}
	ParamError   = ErrorCode{Code: 10004, Msg: "参数错误"}
	BadRequest   = ErrorCode{Code: 10004, Msg: "参数错误"}

	UserNotFound = ErrorCode{Code: 20001, Msg: "用户不存在"}

	AuthorError          = ErrorCode{Code: 20002, Msg: "作者错误"}
	SectionError         = ErrorCode{Code: 20003, Msg: "版块错误"}
	ErrUserAlreadyExist  = ErrorCode{Code: 20004, Msg: "用户已存在"}
	ErrPasswordIncorrect = ErrorCode{Code: 20005, Msg: "密码错误"}
	ErrUserNotExist      = ErrorCode{Code: 20006, Msg: "用户不存在"}
)
