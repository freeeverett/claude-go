package claude_go

import "fmt"

type ResponseError struct {
	Type         string `json:"type"`
	Code         int    `json:"code"`
	ErrorContent struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("code: %d, Type: %s, error_type: %s, message: %s", e.Code, e.Type, e.ErrorContent.Type, e.ErrorContent.Message)
}
