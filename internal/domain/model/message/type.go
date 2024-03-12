package message

type Type string

const (
	Success Type = "success"
	Info    Type = "info"
	Warning Type = "warning"
	Error   Type = "err"
)
