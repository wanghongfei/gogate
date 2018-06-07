package asynclog

import "github.com/alecthomas/log4go"

var Log *AsyncLog

func InitAsyncLog(cfgFile string, queueSize int) {
	Log = NewAsyncLog(cfgFile, queueSize)
}

type AsyncLog struct {
	configFile			string
	logQueue			chan *logData
}

type logData struct {
	Arg0		interface{}
	Args		[]interface{}
	Level		log4go.Level
}

func NewAsyncLog(cfgFile string, queueSize int) *AsyncLog {
	if queueSize < 1 {
		queueSize = 1000
	}

	l := &AsyncLog{
		configFile: cfgFile,
		logQueue: make(chan *logData, queueSize),
	}

	log4go.LoadConfiguration(cfgFile)

	go l.logRoutine()

	return l
}

func (al *AsyncLog) Info(arg0 interface{}, args ...interface{}) {
	data := &logData{
		Arg0: arg0,
		Args: args,
		Level: log4go.INFO,
	}

	al.logQueue <- data
}

func (al *AsyncLog) Debug(arg0 interface{}, args ...interface{}) {
	data := &logData{
		Arg0: arg0,
		Args: args,
		Level: log4go.DEBUG,
	}

	al.logQueue <- data
}

func (al *AsyncLog) Warn(arg0 interface{}, args ...interface{}) {
	data := &logData{
		Arg0: arg0,
		Args: args,
		Level: log4go.WARNING,
	}

	al.logQueue <- data
}

func (al *AsyncLog) Error(arg0 interface{}, args ...interface{}) {
	data := &logData{
		Arg0: arg0,
		Args: args,
		Level: log4go.ERROR,
	}

	al.logQueue <- data
}

func (al *AsyncLog) logRoutine() {
	for record := range al.logQueue {
		if nil == record {
			log4go.Info("log routine exits")
			return
		}

		switch record.Level {
		case log4go.INFO:
			if nil == record.Args {
				log4go.Info(record.Arg0)
			} else {
				log4go.Info(record.Arg0, record.Args...)
			}

		case log4go.DEBUG:
			if nil == record.Args {
				log4go.Debug(record.Arg0)
			} else {
				log4go.Debug(record.Arg0, record.Args...)
			}

		case log4go.WARNING:
			if nil == record.Args {
				log4go.Warn(record.Arg0)
			} else {
				log4go.Warn(record.Arg0, record.Args...)
			}

		case log4go.ERROR:
			if nil == record.Args {
				log4go.Error(record.Arg0)
			} else {
				log4go.Error(record.Arg0, record.Args...)
			}
		}
	}
}

func Info(arg0 interface{}, args ...interface{}) {
	Log.Info(arg0, args...)
}

func Debug(arg0 interface{}, args ...interface{}) {
	Log.Debug(arg0, args...)
}

func Warn(arg0 interface{}, args ...interface{}) {
	Log.Warn(arg0, args...)
}

func Error(arg0 interface{}, args ...interface{}) {
	Log.Error(arg0, args...)
}

