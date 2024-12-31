package helper

type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error interface{} `json:"error,omitempty"`
}

func NewResponse(data interface{}, error error) *Response {

	if error != nil {
		return &Response{
			Data:  data,
			Error: error.Error(),
		}
	}

	return &Response{
		data,
		error,
	}
}
