package error

import (
	"encoding/json"
)

type TError struct {
	Code int32  `json:"code"`
	Desc string `json:"desc"`
}

func (e *TError) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

var (
	ErrServer = &TError{Code: 10000, Desc: "server internal error"}

	ErrRequest = &TError{Code: 20000, Desc: "request params is incorrect"}
)
