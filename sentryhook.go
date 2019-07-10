package sentryhook

import (
	"reflect"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

var levelMap = map[logrus.Level]sentry.Level{
	logrus.TraceLevel: sentry.LevelDebug,
	logrus.DebugLevel: sentry.LevelDebug,
	logrus.InfoLevel:  sentry.LevelInfo,
	logrus.WarnLevel:  sentry.LevelWarning,
	logrus.ErrorLevel: sentry.LevelError,
	logrus.FatalLevel: sentry.LevelFatal,
	logrus.PanicLevel: sentry.LevelFatal,
}

// SentryHook implements logrus.Hook to send errors to sentry.
type SentryHook struct {
	LogLevels []logrus.Level
}

// Levels returns the levels this hook is enabled for. This is a part
// of logrus.Hook.
func (h SentryHook) Levels() []logrus.Level {
	return h.LogLevels
}

// Fire is an event handler for logrus. This is a part of logrus.Hook.
func (h SentryHook) Fire(e *logrus.Entry) error {
	event := sentry.NewEvent()
	event.Level = levelMap[e.Level]
	event.Message = e.Message

	for k, v := range e.Data {
		if k == logrus.ErrorKey {
			event.Exception = []sentry.Exception{
				exceptionFromError(v.(error)),
			}
			continue
		}
		event.Extra[k] = v
	}

	sentry.CaptureEvent(event)

	return nil
}

func exceptionFromError(err error) sentry.Exception {
	stacktrace := sentry.ExtractStacktrace(err)
	if stacktrace == nil {
		stacktrace = sentry.NewStacktrace()
	}
	return sentry.Exception{
		Value:      err.Error(),
		Type:       reflect.TypeOf(err).String(),
		Stacktrace: stacktrace,
	}
}

// New returns a SentryHook with default log levels.
func New() SentryHook {
	return SentryHook{
		LogLevels: []logrus.Level{
			logrus.WarnLevel,
			logrus.ErrorLevel,
			logrus.FatalLevel,
			logrus.PanicLevel,
		},
	}
}
