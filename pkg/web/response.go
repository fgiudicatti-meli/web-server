package web

type Response struct {
	StatusCode int    `json:"code"`
	Data       any    `json:"data"`
	ErrorMsg   string `json:"error"`
}

func NewResponse(code int, data any, err string) Response {
	if code < 300 {
		return Response{code, data, ""}
	}

	return Response{code, nil, err}
}
