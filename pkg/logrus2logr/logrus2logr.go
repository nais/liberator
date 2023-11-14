package logrus2logr

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
)

type Logrus2Logr struct {
	Logger logrus.FieldLogger
	name   string
}

const (
	LevelDebug = 1
	LevelTrace = 2
)

func (l *Logrus2Logr) Init(_ logr.RuntimeInfo) {
}

func (l *Logrus2Logr) Enabled(_ int) bool {
	return l.Logger != nil
}

func (l *Logrus2Logr) Info(level int, msg string, keysAndValues ...interface{}) {
	fields := makeFields(keysAndValues)
	entry := l.Logger.WithFields(fields)
	if level >= LevelTrace {
		entry.Trace(msg)
	} else if level >= LevelDebug {
		entry.Debug(msg)
	} else {
		entry.Info(msg)
	}
}

func (l *Logrus2Logr) Error(err error, msg string, keysAndValues ...interface{}) {
	fields := makeFields(keysAndValues)
	l.Logger.WithFields(fields).Error(fmt.Errorf("%s: %w", msg, err))
}

func (l *Logrus2Logr) WithValues(keysAndValues ...interface{}) logr.LogSink {
	fields := makeFields(keysAndValues)
	return &Logrus2Logr{Logger: l.Logger.WithFields(fields)}
}

func (l *Logrus2Logr) WithName(name string) logr.LogSink {
	if l.name != "" {
		name = fmt.Sprintf("%s.%s", l.name, name)
	}
	return &Logrus2Logr{Logger: l.Logger.WithField("logger_name", name), name: name}
}

func makeFields(keysAndValues []interface{}) logrus.Fields {
	fields := logrus.Fields{}
	for i, value := range keysAndValues {
		if i%2 == 1 {
			key := keysAndValues[i-1].(string)
			fields[key] = value
		}
	}
	return fields
}
