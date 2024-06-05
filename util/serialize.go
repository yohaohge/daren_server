package util

func Pack(code int, msg string, data interface{}) interface{} {
	return map[string]interface{}{
		"code": code,
		"msg":  msg,
		"data": data,
	}
}
