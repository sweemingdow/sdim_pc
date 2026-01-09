package client

import (
	"errors"
	"fmt"
)

type clientHandleError struct {
	err       error // origin error
	reconnect bool  // 是否重连
}

func newClientHandleError(err error, reconnect bool) clientHandleError {
	return clientHandleError{
		err:       err,
		reconnect: reconnect,
	}
}

func (che clientHandleError) Error() string {
	return fmt.Sprintf("reconnect=%t, err=%v", che.reconnect, che.err)
}

func (che clientHandleError) Unwrap() error {
	return che.err
}

func (che clientHandleError) ShouldReconnect() bool {
	return che.reconnect
}

func decodeClientHandleError(err error) (clientHandleError, bool) {
	var che clientHandleError
	if errors.As(err, &che) {
		return che, true
	}

	return che, false
}
