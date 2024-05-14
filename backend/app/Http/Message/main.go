package Message

type Response struct {
	Success bool   `json:"success"`
	Result  string `json:"result"`
	Error   int    `json:"error"`
}
