package http

const Error = "error"
const Data = "data"

type Resp map[string]interface{}

func Unauthorized(payload ...interface{}) Resp {
	return error("unauthorized", payload...)
}

func BadRequest(payload ...interface{}) Resp {
	return error("invalid data", payload...)
}

func Forbidden(payload ...interface{}) Resp {
	return error("forbidden", payload...)
}

func InternalServErr(payload ...interface{}) Resp {
	return error("internal server error", payload...)
}

func Success(payload interface{}) Resp {
	return Resp{"data": payload}
}

func error(def string, payload ...interface{}) Resp {
	if len(payload) == 0 {
		return Resp{
			Error: def,
		}
	}
	return Resp{
		Error: payload,
	}
}
