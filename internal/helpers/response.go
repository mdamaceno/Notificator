package helpers

const (
	GENERIC_MESSAGE       string = "something_went_wrong"
	INVALID_REQUEST       string = "invalid_request_body"
	INTERNAL_SERVER_ERROR string = "internal_server_error"
	NOT_ENOUGH_RECEIVERS  string = "not_enough_receivers"
)

type APIResponse struct {
	Result interface{}       `json:"result"`
	Error  *APIErrorResponse `json:"error"`
}

type APIErrorResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewAPIResponse(result interface{}) APIResponse {
	return APIResponse{
		Result: result,
	}
}

func NewAPIErrorResponse(message string, data interface{}) APIResponse {
	return APIResponse{
		Result: nil,
		Error: &APIErrorResponse{
			Message: message,
			Data:    data,
		},
	}
}
