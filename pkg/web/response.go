package web

type Response struct {
	Code  int    `json:"code"`
	Data  any    `json:"data"`
	Error string `json:"error"`
}

func NewResponse(code int, data any, err string) Response {
	if code < 300 {
		return Response{code, data, ""}
	}

	return Response{code, nil, err}
}
