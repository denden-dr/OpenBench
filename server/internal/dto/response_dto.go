package dto

type APIResponse struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data,omitempty"`
}

func NewAPIResponse(message string, status int, data interface{}) APIResponse {
	return APIResponse{
		Message: message,
		Status:  status,
		Data:    data,
	}
}
