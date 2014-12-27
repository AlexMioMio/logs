// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/issue9/logs/writer"
)

// 预定义的6个log.Logger实例。
var (
	INFO     *log.Logger
	WARN     *log.Logger
	ERROR    *log.Logger
	DEBUG    *log.Logger
	TRACE    *log.Logger
	CRITICAL *log.Logger
)

// 保存INFO,WARN等6个预定义log.Logger实例的io.Writer接口实例，
// 方便在关闭日志时，输出其中缓存的内容。
var conts = writer.NewContainer()

// 从一个xml文件中初始化日志系统。
func InitFromXmlFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()

	return Init(f)
}

// 从一个xml字符串初始化日志系统。
func InitFromXml(xml string) error {
	return Init(bytes.NewBufferString(xml))
}

func init() {
	// 默认初始化为所有的日志实例。
	discardLog := log.New(ioutil.Discard, "", log.LstdFlags)
	INFO = discardLog
	WARN = discardLog
	DEBUG = discardLog
	ERROR = discardLog
	TRACE = discardLog
	CRITICAL = discardLog
}

// 从一个io.Reader初始化日志系统。
// r为一个有xml格式内容的io.Reader实例。
func Init(r io.Reader) error {
	cfg, err := loadFromXml(r)
	if err != nil {
		return err
	}

	if cfg.name != "logs" {
		return fmt.Errorf("顶级元素必须为logs，当前名称为[%v]", cfg.name)
	}

	if len(cfg.attrs) > 0 {
		return fmt.Errorf("logs元素不存在任何属性")
	}

	for name, c := range cfg.items {
		writer, err := c.toWriter()
		if err != nil {
			return err
		}

		w, ok := writer.(*logWriter)
		if !ok {
			return errors.New("二级元素必须为logWriter类型")
		}
		switch strings.ToLower(name) {
		case "info":
			INFO = w.toLogger()
		case "warn":
			WARN = w.toLogger()
		case "debug":
			DEBUG = w.toLogger()
		case "error":
			ERROR = w.toLogger()
		case "trace":
			TRACE = w.toLogger()
		case "critical":
			CRITICAL = w.toLogger()
		default:
			return fmt.Errorf("未知的二级元素:[%v]", name)
		}
		conts.Add(w.c)
	}

	return nil
}

// 输出所有的缓存内容。
func Flush() {
	conts.Flush()
}

// Info相当于INFO.Println(v...)的简写方式
func Info(v ...interface{}) {
	INFO.Println(v...)
}

// Infof相当于INFO.Printf(format, v...)的简写方式
func Infof(format string, v ...interface{}) {
	INFO.Printf(format, v...)
}

// Debug相当于DEBUG.Println(v...)的简写方式
func Debug(v ...interface{}) {
	DEBUG.Println(v...)
}

// Debugf相当于DEBUG.Printf(format, v...)的简写方式
func Debugf(format string, v ...interface{}) {
	DEBUG.Printf(format, v...)
}

// Trace相当于TRACE.Println(v...)的简写方式
func Trace(v ...interface{}) {
	TRACE.Println(v...)
}

// Tracef相当于TRACE.Printf(format, v...)的简写方式
func Tracef(format string, v ...interface{}) {
	TRACE.Printf(format, v...)
}

// Warn相当于WARN.Println(v...)的简写方式
func Warn(v ...interface{}) {
	WARN.Println(v...)
}

// Warnf相当于WARN.Printf(format, v...)的简写方式
func Warnf(format string, v ...interface{}) {
	WARN.Printf(format, v...)
}

// Error相当于ERROR.Println(v...)的简写方式
func Error(v ...interface{}) {
	WARN.Println(v...)
}

// Errorf相当于ERROR.Printf(format, v...)的简写方式
func Errorf(format string, v ...interface{}) {
	WARN.Printf(format, v...)
}

// Critical相当于CRITICAL.Println(v...)的简写方式
func Critical(v ...interface{}) {
	CRITICAL.Println(v...)
}

// Criticalf相当于CRITICAL.Printf(format, v...)的简写方式
func Criticalf(format string, v ...interface{}) {
	CRITICAL.Printf(format, v...)
}

// 向所有的日志输出内容。
func All(v ...interface{}) {
	Info(v...)
	Debug(v...)
	Trace(v...)
	Warn(v...)
	Error(v...)
	Critical(v...)
}

// 向所有的日志输出内容。
func Allf(format string, v ...interface{}) {
	Infof(format, v...)
	Debugf(format, v...)
	Tracef(format, v...)
	Warnf(format, v...)
	Errorf(format, v...)
	Criticalf(format, v...)
}
