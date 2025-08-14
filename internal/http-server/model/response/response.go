package response

type Response struct {
	Status string `json:"status"`
	Err    string `json:"error,omitempty"`
}

func NewStatusOK() Response {
	return Response{
		Status: "OK",
	}
}

func NewStatusError(err string) Response {
	return Response{
		Status: "Error",
		Err:    err,
	}
}
