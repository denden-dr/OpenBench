package dto

type ProblemDetails struct {
	Type          string            `json:"type"`
	Title         string            `json:"title"`
	Status        int               `json:"status"`
	Detail        string            `json:"detail"`
	Instance      string            `json:"instance"`
	InvalidParams map[string]string `json:"invalid_params,omitempty"`
}
