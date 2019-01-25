package log

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"github.com/tinklabs/golibs/cmd"
	"github.com/tinklabs/golibs/utils"
)

var (
	logger *logrus.Logger
)

// Fields wraps logrus.Fields, which is a map[string]interface{}
type Fields logrus.Fields

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Init() {
	logger = logrus.New()
	cf := cmd.GetCmdFlag()

	logFileName := fmt.Sprintf("%s.log", cf.ServerName)
	f, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	// Log as JSON instead of the default ASCII formatter.
	logger.SetFormatter(&logrus.JSONFormatter{})
	// logger.SetFormatter(&logrus.TextFormatter{
	// 	DisableColors: false,
	// 	FullTimestamp: true,
	// })

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example

	// Only log the warning severity or above.
	if cmd.IsDebug() {
		logger.SetLevel(logrus.DebugLevel)
		logger.SetOutput(os.Stdout)
	} else {
		logger.SetLevel(logrus.InfoLevel)
		mw := io.MultiWriter(os.Stdout, f)
		logger.SetOutput(mw)
	}
}

// LoggerWithWriter instance a Logger middleware with the specified writter buffer.
// Example: os.Stdout, a file opened in write mode, a socket...
func Logger() gin.HandlerFunc {

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path

		buf, _ := ioutil.ReadAll(c.Request.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))

		c.Request.Body = rdr2

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		reqStr := utils.ReadBody(rdr1)
		if c.ContentType() == "multipart/form-data" {
			f, _ := c.FormFile("file")
			if f != nil {
				reqStr = f.Filename
			}
		}
		DebugWithFields("", Fields{"request-id": c.GetHeader("Request-Id"), "request": reqStr, "response": blw.body.String(), "clientIP": clientIP,
			"path": path, "method": method, "statusCode": statusCode, "latency": latency})

	}
}

func SetLogFormatter(formatter logrus.Formatter) {
	logger.Formatter = formatter
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Debug(args...)
	}
}

// Debug logs a message with fields at level Debug on the standard logger.
func DebugWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Debug(l)
	}
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Info(args...)
	}
}

// Debug logs a message with fields at level Debug on the standard logger.
func InfoWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Info(l)
	}
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Warn(args...)
	}
}

// Debug logs a message with fields at level Debug on the standard logger.
func WarnWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.WarnLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Warn(l)
	}
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Error(args...)
	}
}

// Debug logs a message with fields at level Debug on the standard logger.
func ErrorWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.ErrorLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Error(l)
	}
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(args ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Fatal(args...)
	}
}

// Debug logs a message with fields at level Debug on the standard logger.
func FatalWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.FatalLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Fatal(l)
	}
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	if logger.Level >= logrus.PanicLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Panic(args...)
	}
}

// Debug logs a message with fields at level Debug on the standard logger.
func PanicWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.PanicLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Panic(l)
	}
}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
