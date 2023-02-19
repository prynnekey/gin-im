package response

func Response(code int, data interface{}, msg string) map[string]interface{} {
	return map[string]interface{}{
		"code": code,
		"data": data,
		"msg":  msg,
	}
}

func Success(data interface{}, msg string) map[string]interface{} {
	return Response(200, data, msg)
}

func Fail(data interface{}, msg string) map[string]interface{} {
	return Response(500, data, msg)
}
