package time

import (
	log "github.com/sirupsen/logrus"
	"time"
)

const timeCost = "[Time-Cost]: "

type InfoLogger interface {
	Info(v ...interface{})
}

type PrintlnLogger interface {
	Println(v ...interface{})
}

// StdLogger define a default log instance which implement Logger
type StdLogger struct{}

func (l StdLogger) Println(v ...interface{}) {
	log.Println(v...)
}

// TimeCost
// If you want to measure the execution time of a method
// and find it tedious and ugly to write `time.Now()` at both ends of the method
// you can consider using this method `TimeCost` in defer
// besides,it provider you to offer a prefix(such as `methodName`:5s) and a log instance which implement the InfoLogger / PrintlnLogger
// and it will use the instance.Println/instance.Info to output the time cost
func TimeCost(start time.Time, logInstance interface{}, prefix string) {
	if logInstance == nil {
		logInstance = StdLogger{}
	}
	// if logInstance is nil,use `log.Println()`
	if logInstance == nil {
		StdLogger{}.Println(timeCost, prefix, time.Since(start))
		return
	}

	if l, ok := logInstance.(InfoLogger); ok {
		l.Info(timeCost, prefix, time.Since(start))
		return
	}

	if l, ok := logInstance.(PrintlnLogger); ok {
		l.Println(timeCost, prefix, time.Since(start))
		return
	}
	return
}
