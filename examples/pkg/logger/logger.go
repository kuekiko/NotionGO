package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
)

// Init 初始化日志
func Init() error {
	// 创建日志目录
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 创建日志文件
	logFile := filepath.Join(logDir, fmt.Sprintf("notion_%s.log", time.Now().Format("2006-01-02")))
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("创建日志文件失败: %v", err)
	}

	// 初始化日志记录器
	infoLogger = log.New(file, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(file, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	debugLogger = log.New(file, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}

// Info 记录信息日志
func Info(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	fmt.Println(msg)
	if infoLogger != nil {
		infoLogger.Output(2, msg)
	}
}

// Error 记录错误日志
func Error(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	fmt.Println(msg)
	if errorLogger != nil {
		errorLogger.Output(2, msg)
	}
}

// Debug 记录调试日志
func Debug(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	if debugLogger != nil {
		debugLogger.Output(2, msg)
	}
}
