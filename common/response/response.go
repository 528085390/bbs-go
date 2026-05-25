package response

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func Ok(data interface{}) Response {
	return Response{Code: 0, Msg: "success", Data: data}
}
func Error(code int, message string, data interface{}) Response {
	return Response{Code: code, Msg: message, Data: data}
}
