package httperror

import (
	"context"
	"fmt"
)

type HttpError struct {
	message    string
	ctx        context.Context
	statusCode int
}

func (e HttpError) Error() string {
	return e.message
}

func (e HttpError) Context() context.Context {
	return e.ctx
}

func (e HttpError) Status() int {
	return e.statusCode
}

func New(ctx context.Context, httpStatus int, format string, a ...any) HttpError {
	return HttpError{ctx: ctx, statusCode: httpStatus, message: fmt.Errorf(format, a...).Error()}
}
