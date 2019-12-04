package tracker

import (
	"context"
	"io"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/sqmmm/finance-app/internal/logger"
)

const (
	requestIDName = "request_id"
)

type tracker struct {
	logger    *logrus.Logger
	requestID int64
}

var log *logrus.Logger

func InitLog(out io.Writer) error {
	log = &logrus.Logger{
		Out:       out,
		Formatter: new(logrus.TextFormatter),
		Level:     logrus.DebugLevel,
	}

	return nil
}

type manager struct{}

func NewTrackerManager() *manager {
	return &manager{}
}

func (m *manager) Log() logger.Logger {
	return &tracker{logger: log}
}

func (m *manager) LogCtx(ctx context.Context) logger.Logger {
	//if value does not exist tracker.go.id=0 - ok! (warning)
	requestID, okR := ctx.Value(requestIDName).(int64)
	tr := &tracker{logger: log, requestID: requestID}

	if !okR || requestID == 0 {
		tr.Warningf("requestID for tracker.go does not exist in context")
	}

	return tr
}

func (m *manager) SetRequestID(ctx context.Context, requestID int64) context.Context {
	ctx = context.WithValue(ctx, requestIDName, requestID)

	return ctx
}

func (m *manager) GetRequestID(ctx context.Context) (int64, error) {
	requestID, ok := ctx.Value(requestIDName).(int64)
	if !ok {
		return 0, errors.New("requestID does not exist in context")
	}

	return requestID, nil
}

func (t *tracker) Errorf(msg string, args ...interface{}) {
	t.logger.WithField(requestIDName, t.requestID).Errorf(msg, args...)
}

func (t *tracker) Infof(msg string, args ...interface{}) {
	t.logger.WithField(requestIDName, t.requestID).Infof(msg, args...)
}

func (t *tracker) Debugf(msg string, args ...interface{}) {
	t.logger.WithField(requestIDName, t.requestID).Debugf(msg, args...)
}

func (t *tracker) Warningf(msg string, args ...interface{}) {
	t.logger.WithField(requestIDName, t.requestID).Warningf(msg, args...)
}

func (t *tracker) Writer() io.Writer {
	return t.logger.Writer()
}
