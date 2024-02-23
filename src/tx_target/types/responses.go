package types

type SendRawTransactionResponse struct {
	Result string         `json:"result"`
	Error  *ResponseError `json:"error,omitempty"`
}

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
