package common

type ResultUtil struct {
}

var ResultUtils = &ResultUtil{}

func (r *ResultUtil) Success(data interface{}) *ResultCommon {
	rc := ResultCommon{Code: "200"}
	rc.Data = data
	return &rc
}

func (r *ResultUtil) Fail(code, message string) *ResultCommon {
	return &ResultCommon{Code: code, Message: message}
}

type ResultCommon struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
