package exactonline

type ExactOnlineError struct {
	Err ExactOnlineInnerError `json:"error"`
}
type ExactOnlineInnerError struct {
	Code    string                  `json:"code"`
	Message ExactOnlineErrorMessage `json:"message"`
}

type ExactOnlineErrorMessage struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}
