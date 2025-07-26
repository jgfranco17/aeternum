package httperror

import (
	"context"
	"fmt"
)

type InputError struct {
	message string
	ctx     context.Context
}

func (e InputError) Error() string {
	return e.message
}

func (e InputError) Context() context.Context {
	return e.ctx
}

func NewInputError(ctx context.Context, format string, a ...any) InputError {
	return InputError{ctx: ctx, message: fmt.Errorf(format, a...).Error()}
}
