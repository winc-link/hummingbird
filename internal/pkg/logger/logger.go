/*******************************************************************************
 * Copyright 2019 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

/*
Package logger provides a client for integration with the support-logging service. The client can also be configured
to write logs to a local file rather than sending them to a service.
*/
package logger

// Logging client for the Go implementation of edgexfoundry

import (
	"fmt"
	stdLog "log"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Colors
const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

// LoggingClient defines the interface for logging operations.
type LoggingClient interface {
	// SetLogLevel sets minimum severity log level. If a logging method is called with a lower level of severity than
	// what is set, it will result in no output.
	SetLogLevel(logLevel string) error
	// LogLevel returns the current log level setting
	LogLevel() string
	// Debug logs a message at the DEBUG severity level
	Debug(msg string, args ...interface{})
	// Error logs a message at the ERROR severity level
	Error(msg string, args ...interface{})
	// Info logs a message at the INFO severity level
	Info(msg string, args ...interface{})
	// Trace logs a message at the TRACE severity level
	Trace(msg string, args ...interface{})
	// Warn logs a message at the WARN severity level
	Warn(msg string, args ...interface{})
	// Debugf logs a formatted message at the DEBUG severity level
	Debugf(msg string, args ...interface{})
	// Errorf logs a formatted message at the ERROR severity level
	Errorf(msg string, args ...interface{})
	// Infof logs a formatted message at the INFO severity level
	Infof(msg string, args ...interface{})
	// Tracef logs a formatted message at the TRACE severity level
	Tracef(msg string, args ...interface{})
	// Warnf logs a formatted message at the WARN severity level
	Warnf(msg string, args ...interface{})
}

type edgeXLogger struct {
	owningServiceName string
	level             zap.AtomicLevel
	rootLogger        *zap.Logger
	levelLoggers      map[string]*zap.Logger
}

const (
	LogPathEnvName  = "LOG_PATH"
	LogLevel        = "LOG_LEVEL"
	DefaultLogLevel = InfoLog
	DefaultLogPath  = ""
)

var (
	zapLevels = map[string]zapcore.Level{
		DebugLog: zap.DebugLevel,
		InfoLog:  zap.InfoLevel,
		WarnLog:  zap.WarnLevel,
		ErrorLog: zap.ErrorLevel,
	}

	levels = map[zapcore.Level]string{
		zap.DebugLevel: DebugLog,
		zap.InfoLevel:  InfoLog,
		zap.WarnLevel:  WarnLog,
		zap.ErrorLevel: ErrorLog,
	}
)

type LogMessage struct {
	Time        string        `json:"time"`
	ServiceName string        `json:"service_name"`
	Caller      string        `json:"caller"`
	Message     []interface{} `json:"message"`
}

func NewZapLogger(atomLevel zap.AtomicLevel, logPath string) (zapLog *zap.Logger, err error) {
	stdLog.SetFlags(stdLog.LstdFlags | stdLog.Llongfile)
	// 选择自定义日志样式
	encoderConfig := zapcore.EncoderConfig{
		MessageKey: "msg",
		//StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if logPath != "" {
		logDir := path.Dir(logPath)
		_, fileErr := os.Stat(logDir)
		if fileErr != nil || !os.IsExist(fileErr) {
			err = os.MkdirAll(logDir, os.ModePerm)
			if err != nil {
				stdLog.Fatal(fmt.Sprintf("ERROR mkdir dir %s err %+v ", logDir, err))
			}
		}
		// 打印到文件，自动分裂
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    10, // megabytes
			MaxBackups: 3,
			MaxAge:     7, // days
		})
		core := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			w,
			atomLevel,
		)
		zapLog = zap.New(core, zap.AddCaller())
	} else {
		// 打印到控制台
		cfg := zap.NewProductionConfig()
		cfg.Level = atomLevel
		cfg.Encoding = "console"
		cfg.EncoderConfig = encoderConfig
		zapLog, err = cfg.Build()
		if err != nil {
			stdLog.Fatal("ERROR ", err)
			return
		}
	}
	return
}

// NewClient creates an instance of LoggingClient
func NewClient(owningServiceName string, logLevel, logPath string) LoggingClient {
	if logLevel == "" {
		// 从Env环境变量获取日志等级
		logLevel = os.Getenv(LogLevel)
		if logLevel == "" {
			logLevel = DefaultLogLevel
		}
	}
	if logPath == "" {
		logPath = os.Getenv(LogLevel)
	}

	if !isValidLogLevel(logLevel) {
		logLevel = DebugLog
	}

	// Set up logging client
	lc := edgeXLogger{
		owningServiceName: owningServiceName,
		level:             zap.NewAtomicLevelAt(zapLevels[logLevel]),
	}

	var err error
	lc.rootLogger, err = NewZapLogger(lc.level, logPath)
	if err != nil {
		return nil
	}

	lc.levelLoggers = make(map[string]*zap.Logger)
	for _, level := range logLevels() {
		lc.levelLoggers[level] = lc.rootLogger
	}

	return lc
}

// LogLevels returns an array of the possible log levels in order from most to least verbose.
func logLevels() []string {
	return []string{
		TraceLog,
		DebugLog,
		InfoLog,
		WarnLog,
		ErrorLog,
	}
}

func isValidLogLevel(l string) bool {
	for _, name := range logLevels() {
		if name == l {
			return true
		}
	}
	return false
}

func (lc edgeXLogger) check(logLevel string) bool {
	l := zapLevels[logLevel]
	return lc.level.Enabled(l)
}

func (lc edgeXLogger) log(logLevel string, formatted bool, msg string, args ...interface{}) {
	// Check minimum log level
	if lc.check(logLevel) == false {
		return
	}

	if args == nil {
		args = []interface{}{msg}
	} else if formatted {
		args = []interface{}{fmt.Sprintf(msg, args...)}
	} else {
		if len(msg) > 0 {
			args = append([]interface{}{msg}, args...)
		}
	}

	argData := make([]string, 0)
	for _, arg := range args {
		argData = append(argData, fmt.Sprintf("%+v", arg))
	}
	msg = strings.Join(argData, ",")

	//_, file, line, _ := runtime.Caller(2)
	//idx := strings.LastIndexByte(file, '/')
	//caller := file[idx+1:] + ":" + strconv.Itoa(line)

	// 日志样式设置
	pc, file, line, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	paths := strings.Split(funcName, "/")
	pack := strings.Split(funcName[strings.LastIndexByte(funcName, '/')+1:], ".")[0]
	paths = append(paths[:len(paths)-1], pack, strings.Split(file[strings.LastIndexByte(file, '/')+1:], ".")[0])
	funcPath := strings.Join(paths, "/")
	caller := funcPath + ".go:" + strconv.Itoa(line) + " "
	now := time.Now().Format("2006/01/02 15:04:05")
	debugStr := Yellow + "[debug] " + Reset
	debugCaller := caller
	infoStr := Green + "[info] " + Reset
	infoCaller := caller
	warnStr := Magenta + "[warn] " + Reset
	warnCaller := caller
	errStr := Red + "[error] " + Reset
	errCaller := caller
	var message = fmt.Sprintf("%v [%v] ", now, lc.owningServiceName)

	// 日志输出
	switch logLevel {
	case DebugLog:
		lc.levelLoggers[logLevel].Debug(debugStr + message + debugCaller + msg)
	case InfoLog:
		lc.levelLoggers[logLevel].Info(infoStr + message + infoCaller + msg)
	case WarnLog:
		lc.levelLoggers[logLevel].Warn(warnStr + message + warnCaller + msg)
	case ErrorLog:
		lc.levelLoggers[logLevel].Error(errStr + message + errCaller + msg)
	}
}

func (lc edgeXLogger) SetLogLevel(logLevel string) error {
	lc.level.SetLevel(zapLevels[logLevel])
	return nil
}

func (lc edgeXLogger) LogLevel() string {
	l := lc.level.Level()
	return levels[l]
}

func (lc edgeXLogger) Info(msg string, args ...interface{}) {
	lc.log(InfoLog, false, msg, args...)
}

func (lc edgeXLogger) Trace(msg string, args ...interface{}) {
	lc.log(TraceLog, false, msg, args...)
}

func (lc edgeXLogger) Debug(msg string, args ...interface{}) {
	lc.log(DebugLog, false, msg, args...)
}

func (lc edgeXLogger) Warn(msg string, args ...interface{}) {
	lc.log(WarnLog, false, msg, args...)
}

func (lc edgeXLogger) Error(msg string, args ...interface{}) {
	lc.log(ErrorLog, false, msg, args...)
}

func (lc edgeXLogger) Infof(msg string, args ...interface{}) {
	lc.log(InfoLog, true, msg, args...)
}

func (lc edgeXLogger) Tracef(msg string, args ...interface{}) {
	lc.log(TraceLog, true, msg, args...)
}

func (lc edgeXLogger) Debugf(msg string, args ...interface{}) {
	lc.log(DebugLog, true, msg, args...)
}

func (lc edgeXLogger) Warnf(msg string, args ...interface{}) {
	lc.log(WarnLog, true, msg, args...)
}

func (lc edgeXLogger) Errorf(msg string, args ...interface{}) {
	lc.log(ErrorLog, true, msg, args...)
}

// Build the log entry object
func (lc edgeXLogger) buildLogEntry(logLevel string, msg string, args ...interface{}) LogEntry {
	res := LogEntry{}
	res.Level = logLevel
	res.Message = msg
	res.Args = args
	res.OriginService = lc.owningServiceName

	return res
}
