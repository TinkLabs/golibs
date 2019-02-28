package error

import (
	"encoding/json"
	"fmt"
)

type TError struct {
	Code  int    `json:"code"`
	Desc  string `json:"desc"`
	Extra string `json:"extra,omitempty"`
}

func (e *TError) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

var (
	ErrServer = &TError{Code: 10000, Desc: "server internal error"}
	ErrConsul = &TError{Code: 10001, Desc: "consul error"}

	ErrRequest = &TError{Code: 20000, Desc: "request params is incorrect"}
)

func (e *TError) AddExtra(extra string) (err *TError) {
	temp := *e

	err = &temp
	err.Extra = extra

	return err
}

func (e *TError) Message() (msg string) {
	if e.Extra != "" {
		msg = fmt.Sprintf("%s(%s)", e.Desc, e.Extra)
	} else {
		msg = e.Desc
	}

	return msg
}
