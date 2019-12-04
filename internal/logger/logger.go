package logger

import (
	"context"
	"io"
)

type LoggerManager interface {
	LogCtx(ctx context.Context) Logger
	SetRequestID(ctx context.Context, requestID int64) context.Context
	Log() Logger
}

type Logger interface {
	Errorf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Debugf(msg string, args ...interface{})
	Warningf(msg string, args ...interface{})
	Writer() io.Writer
}
