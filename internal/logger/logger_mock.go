package logger

import (
	"context"
	"io"
)

type LoggerManagerMock struct{}
type loggerMock struct{}

func (m *LoggerManagerMock) Log() Logger {
	return &loggerMock{}
}

func (m *LoggerManagerMock) LogCtx(ctx context.Context) Logger {
	return &loggerMock{}
}

func (m *LoggerManagerMock) SetDeviceID(ctx context.Context, id string) context.Context {
	return ctx
}

func (m *LoggerManagerMock) SetRequestID(ctx context.Context, id int64) context.Context {
	return ctx
}

func (t *loggerMock) Errorf(msg string, args ...interface{}) {}

func (t *loggerMock) Infof(msg string, args ...interface{}) {}

func (t *loggerMock) Debugf(msg string, args ...interface{}) {}

func (t *loggerMock) Warningf(msg string, args ...interface{}) {}

func (t *loggerMock) Writer() io.Writer {
	return nil
}
