package serializer

type Response struct {
	Code int    `json:"code"`
	Data any    `json:"data,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

const (
	CodeCliParErr = 401 //客户端传入参数有问题
	CodeSerErr    = 500 //服务器出现问题通用报错
)

func Err(code int, msg string, err error) Response {
	if err != nil {
		msg += err.Error()
	}
	return Response{
		Code: code,
		Msg:  msg,
	}
}

func CliParErr(msg string, err error) Response {
	if msg == "" {
		msg = "传入参数有误！"
	}
	return Err(CodeCliParErr, msg, err)
}

func SerErr(msg string, err error) Response {
	if msg == "" {
		msg = "服务器未知错误！"
	}
	return Err(CodeSerErr, msg, err)
}
