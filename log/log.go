package log

import (
	stdLog "log"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var logrusLogger *logrus.Logger
var logrusEntry *logrus.Entry
var stdLogLogger *stdLog.Logger

func init() {
	logrusLogger = logrus.New()
	setDefaultLogLevel()

	logrusLogger.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: false,
	})

	logrusEntry = logrusLogger.WithField("module", "mongoid")
}

//
func setDefaultLogLevel() {
	defaultLevel := logrus.InfoLevel
	switch strings.ToLower(os.Getenv("MONGOID_LOG_LEVEL")) {
	case "error":
		logrusLogger.SetLevel(logrus.ErrorLevel)
	case "warn":
		logrusLogger.SetLevel(logrus.WarnLevel)
	case "info":
		logrusLogger.SetLevel(logrus.InfoLevel)
	case "debug":
		logrusLogger.SetLevel(logrus.DebugLevel)
	case "trace":
		logrusLogger.SetLevel(logrus.TraceLevel)
	default:
		logrusLogger.SetLevel(defaultLevel)
	}
}

// func SetLogger(logger interface{})

func fieldLogger() logrus.Ext1FieldLogger {
	switch {
	case logrusEntry != nil:
		return logrusEntry
	case logrusLogger != nil:
		return logrusLogger
	}
	return nil
}

func stdLogger() logrus.StdLogger {
	if stdLogLogger != nil {
		return stdLogLogger
	}
	return nil
}

// Fatal methods //////////////////////////////////////////////////////////////

// Fatal ...
func Fatal(v ...interface{}) {
	if logger := fieldLogger(); logger != nil {
		logger.Fatal(v...)
	} else if logger := stdLogger(); logger != nil {
		logger.Fatal(v...)
	} else {
		stdLog.Fatal(v...)
	}
}

// Fatalf ...
func Fatalf(format string, v ...interface{}) {
	if logger := fieldLogger(); logger != nil {
		logger.Fatalf(format, v...)
	} else if logger := stdLogger(); logger != nil {
		logger.Fatalf(format, v...)
	} else {
		stdLog.Fatalf(format, v...)
	}
}

// Fatalln ...
func Fatalln(v ...interface{}) {
	if logger := fieldLogger(); logger != nil {
		logger.Fatalln(v...)
	} else if logger := stdLogger(); logger != nil {
		logger.Fatalln(v...)
	} else {
		stdLog.Fatalln(v...)
	}
}

// Panic methods //////////////////////////////////////////////////////////////

// recovers from the panic raised by logrus Panic functions (only)
func recoverLogrusPanic() {
	if err := recover(); err != nil {
		// log.Warn("Recovered")
		// log.Warn(reflect.TypeOf(err))
		switch v := err.(type) {
		case (*logrus.Entry):
		case (logrus.Entry):
		case (*logrus.Logger):
		case (logrus.Logger):
		default:
			panic(v)
		}
	}

}

// Panic ...
func Panic(v ...interface{}) {
	if logger := fieldLogger(); logger != nil {
		func() {
			defer recoverLogrusPanic()
			logger.Panic(v...)
		}()
	} else if logger := stdLogger(); logger != nil {
		func() {
			defer recoverLogrusPanic()
			logger.Panic(v...)
		}()
	} else {
		stdLog.Panic(v...)
	}
	if len(v) == 1 {
		panic(v[0])
	}
	panic(v)
}

// Panicf ...
func Panicf(format string, v ...interface{}) {
	if logger := fieldLogger(); logger != nil {
		logger.Panicf(format, v...)
	} else if logger := stdLogger(); logger != nil {
		logger.Panicf(format, v...)
	} else {
		stdLog.Panicf(format, v...)
	}
}

// Panicln ...
func Panicln(v ...interface{}) {
	if logger := fieldLogger(); logger != nil {
		logger.Panicln(v...)
	} else if logger := stdLogger(); logger != nil {
		logger.Panicln(v...)
	} else {
		stdLog.Panicln(v...)
	}
}

// Print methods //////////////////////////////////////////////////////////////

// Print ...
func Print(v ...interface{}) {
	if logger := fieldLogger(); logger != nil {
		logger.Print(v...)
	} else if logger := stdLogger(); logger != nil {
		logger.Print(v...)
	} else {
		stdLog.Print(v...)
	}
}

// Printf ...
func Printf(format string, v ...interface{}) {
	if logger := fieldLogger(); logger != nil {
		logger.Printf(format, v...)
	} else if logger := stdLogger(); logger != nil {
		logger.Printf(format, v...)
	} else {
		stdLog.Printf(format, v...)
	}
}

// Println ...
func Println(v ...interface{}) {
	if logger := fieldLogger(); logger != nil {
		logger.Println(v...)
	} else if logger := stdLogger(); logger != nil {
		logger.Println(v...)
	} else {
		stdLog.Println(v...)
	}
}

// Syslog methods

// Error methods //////////////////////////////////////////////////////////////

// Error ...
func Error(v ...interface{}) {
	if logger := fieldLogger(); logger != nil {
		logger.Error(v...)
	} else if logger := stdLogger(); logger != nil {
		logger.Print(v...)
	} else {
		stdLog.Print(v...)
	}
}

// Errorf ...
func Errorf(format string, v ...interface{}) {
	if logger := fieldLogger(); logger != nil {
		logger.Errorf(format, v...)
	} else if logger := stdLogger(); logger != nil {
		logger.Printf(format, v...)
	} else {
		stdLog.Printf(format, v...)
	}
}

// Errorln ...
func Errorln(v ...interface{}) {
	if logger := fieldLogger(); logger != nil {
		logger.Errorln(v...)
	} else if logger := stdLogger(); logger != nil {
		logger.Println(v...)
	} else {
		stdLog.Println(v...)
	}
}

// Warn methods ///////////////////////////////////////////////////////////////

// Warn ...
func Warn(v ...interface{}) {
	if logger := fieldLogger(); logger != nil {
		logger.Warn(v...)
	} else if logger := stdLogger(); logger != nil {
		logger.Print(v...)
	} else {
		logger.Print(v...)
	}
}

// Warnf ...
func Warnf(format string, v ...interface{}) {
	if logger := fieldLogger(); logger != nil {
		logger.Warnf(format, v...)
	} else if logger := stdLogger(); logger != nil {
		logger.Printf(format, v...)
	} else {
		logger.Printf(format, v...)
	}
}

// Warnln ...
func Warnln(v ...interface{}) {
	if logger := fieldLogger(); logger != nil {
		logger.Warnln(v...)
	} else if logger := stdLogger(); logger != nil {
		logger.Println(v...)
	} else {
		logger.Println(v...)
	}
}

// Info methods ///////////////////////////////////////////////////////////////

// Info ...
func Info(v ...interface{}) {
	if logger := fieldLogger(); logger != nil {
		logger.Info(v...)
	} else if logger := stdLogger(); logger != nil {
		logger.Print(v...)
	} else {
		logger.Print(v...)
	}
}

// Infof ...
func Infof(format string, v ...interface{}) {
	if logger := fieldLogger(); logger != nil {
		logger.Infof(format, v...)
	} else if logger := stdLogger(); logger != nil {
		logger.Printf(format, v...)
	} else {
		logger.Printf(format, v...)
	}
}

// Infoln ...
func Infoln(v ...interface{}) {
	if logger := fieldLogger(); logger != nil {
		logger.Infoln(v...)
	} else if logger := stdLogger(); logger != nil {
		logger.Println(v...)
	} else {
		logger.Println(v...)
	}
}

// Debug methods //////////////////////////////////////////////////////////////

// Debug ...
func Debug(v ...interface{}) {
	if logger := fieldLogger(); logger != nil {
		logger.Debug(v...)
	} else if logger := stdLogger(); logger != nil {
		logger.Print(v...)
	} else {
		logger.Print(v...)
	}
}

// Debugf ...
func Debugf(format string, v ...interface{}) {
	if logger := fieldLogger(); logger != nil {
		logger.Debugf(format, v...)
	} else if logger := stdLogger(); logger != nil {
		logger.Printf(format, v...)
	} else {
		logger.Printf(format, v...)
	}
}

// Debugln ...
func Debugln(v ...interface{}) {
	if logger := fieldLogger(); logger != nil {
		logger.Debugln(v...)
	} else if logger := stdLogger(); logger != nil {
		logger.Println(v...)
	} else {
		logger.Println(v...)
	}
}

// Trace methods //////////////////////////////////////////////////////////////

// Trace ...
func Trace(v ...interface{}) {
	if logger := fieldLogger(); logger != nil {
		logger.Trace(v...)
	} else if logger := stdLogger(); logger != nil {
		logger.Print(v...)
	} else {
		logger.Print(v...)
	}
}

// Tracef ...
func Tracef(format string, v ...interface{}) {
	if logger := fieldLogger(); logger != nil {
		logger.Tracef(format, v...)
	} else if logger := stdLogger(); logger != nil {
		logger.Printf(format, v...)
	} else {
		logger.Printf(format, v...)
	}
}

// Traceln ...
func Traceln(v ...interface{}) {
	if logger := fieldLogger(); logger != nil {
		logger.Traceln(v...)
	} else if logger := stdLogger(); logger != nil {
		logger.Println(v...)
	} else {
		logger.Println(v...)
	}
}
