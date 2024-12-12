package logging

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Level int

type LogFields map[string]interface{}

var (
	F *os.File

	DefaultPrefix      = ""
	DefaultCallerDepth = 2

	logger     *log.Logger
	logPrefix  = ""
	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

// 启动前先生成日志的处理逻辑
// - 日志文件创建和日志文件设置为写入
func Setup() {
	filePath := getLogFilePath()
	fileName := getLogFileName()
	F, err := openLogFile(fileName, filePath)
	if err != nil {
		log.Fatalln(err)
	}

	// 创建一个日志打印器
	logger = log.New(F, DefaultPrefix, log.LstdFlags)
}

// 设置日志打印的前缀
func setPrefix(level Level) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}

	logger.SetPrefix(logPrefix)
}

func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	logger.Println(v...)
}

func DebugF(format string, v ...any) {
	setPrefix(DEBUG)
	logger.Printf(format, v...)
}

func Info(v ...interface{}) {
	setPrefix(INFO)
	logger.Println(v...)
}

func InfoF(format string, v ...any) {
	setPrefix(INFO)
	logger.Printf(format, v...)
}

func Warn(v ...interface{}) {
	setPrefix(WARN)
	logger.Println(v...)
}

func WarnF(format string, v ...any) {
	setPrefix(WARN)
	logger.Printf(format, v...)
}

func Error(v ...interface{}) {
	setPrefix(ERROR)
	logger.Println(v...)
}

func ErrorF(format string, v ...any) {
	setPrefix(ERROR)
	logger.Printf(format, v...)
}

func Fatal(v ...interface{}) {
	setPrefix(FATAL)
	logger.Fatalln(v...)
}

func FatalF(format string, v ...any) {
	setPrefix(FATAL)
	logger.Printf(format, v...)
}

// 以json的格式输出
func InfoFiled(fields LogFields) {
	setPrefix(INFO)
	data, err := json.Marshal(fields)
	if err != nil {
		ErrorF("InfoFiled Failed error: %v", err)
		return
	}
	logger.Println(string(data))
}
