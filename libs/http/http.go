package http

type Resp map[string]interface{}

func Unauthorized() Resp {
	return Resp{
		"error": "unauthorized",
	}
}

func Success(payload interface{}) Resp {
	return Resp{"data": payload}
}
