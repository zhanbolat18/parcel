package http

const Error = "error"
const Data = "data"

type Resp map[string]interface{}

func Unauthorized(payload ...interface{}) (Resp, error) {
	if len(payload) == 0 {
		return Resp{
			Error: "unauthorized",
		}, nil
	}
	return Resp{
		Error: payload,
	}, nil
}

func Success(payload interface{}) Resp {
	return Resp{"data": payload}
}
