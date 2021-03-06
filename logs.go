// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/issue9/logs/internal/config"
	"github.com/issue9/logs/writers"
)

// 保存 info、warn 等6个预定义 log.Logger 的 io.Writer 接口实例，
// 方便在关闭日志时，输出其中缓存的内容。
var conts = writers.NewContainer()

// 预定义的6个 log.Logger 实例。
var (
	info, warn, erro, debug, trace, critical *log.Logger
)

// 从一个 XML 文件中初始化日志系统。
// 再次调用该函数，将会根据新的配置文件重新初始化日志系统。
func InitFromXMLFile(path string) error {
	cfg, err := config.ParseXMLFile(path)
	if err != nil {
		return err
	}
	return initFromConfig(cfg)
}

// 从一个 XML 字符串初始化日志系统。
// 再次调用该函数，将会根据新的配置文件重新初始化日志系统。
func InitFromXMLString(xml string) error {
	cfg, err := config.ParseXMLString(xml)
	if err != nil {
		return err
	}
	return initFromConfig(cfg)
}

// 从 config.Config 中初始化整个 logs 系统
func initFromConfig(cfg *config.Config) error {
	if conts.Len() > 0 { // 加载新配置文件。先输出旧的内容。
		Flush()
		conts.Clear()

		// 重置为空值
		info = nil
		critical = nil
		debug = nil
		trace = nil
		warn = nil
		erro = nil
	}

	for name, c := range cfg.Items {
		flag := 0
		flagStr, found := c.Attrs["flag"]
		if found && (flagStr != "") {
			flag, found = flagMap[strings.ToLower(flagStr)]
			if !found {
				return fmt.Errorf("未知的Flag参数:[%v]", flagStr)
			}
		}

		cont, err := toWriter(c)
		if err != nil {
			return err
		}
		l := log.New(cont, c.Attrs["prefix"], flag)

		switch name {
		case "info":
			info = l
		case "warn":
			warn = l
		case "debug":
			debug = l
		case "error":
			erro = l
		case "trace":
			trace = l
		case "critical":
			critical = l
		}
		conts.Add(cont)
	}

	return nil
}

// 输出所有的缓存内容。
// 若是通过 os.Exit() 退出程序的，在执行之前，
// 一定记得调用 Flush() 输出可能缓存的日志内容。
func Flush() {
	conts.Flush()
}

// 获取 INFO 级别的 log.Logger 实例，在未指定 info 级别的日志时，该实例返回一个 nil。
func INFO() *log.Logger {
	return info
}

// Info 相当于 INFO().Println(v...) 的简写方式
// Info 函数默认是带换行符的，若需要不带换行符的，请使用 DEBUG().Print() 函数代替。
// 其它相似函数也有类型功能。
func Info(v ...interface{}) {
	if info == nil {
		return
	}

	info.Println(v...)
}

// Infof 相当于 INFO().Printf(format, v...) 的简写方式
func Infof(format string, v ...interface{}) {
	if info == nil {
		return
	}

	info.Printf(format, v...)
}

// 获取 DEBUG 级别的 log.Logger 实例，在未指定 debug 级别的日志时，该实例返回一个 nil。
func DEBUG() *log.Logger {
	return debug
}

// Debug 相当于 DEBUG().Println(v...) 的简写方式
func Debug(v ...interface{}) {
	if debug == nil {
		return
	}

	debug.Println(v...)
}

// Debugf 相当于 DEBUG().Printf(format, v...) 的简写方式
func Debugf(format string, v ...interface{}) {
	if debug == nil {
		return
	}

	debug.Printf(format, v...)
}

// 获取 TRACE 级别的 log.Logger 实例，在未指定 trace 级别的日志时，该实例返回一个 nil。
func TRACE() *log.Logger {
	return trace
}

// Trace 相当于 TRACE().Println(v...) 的简写方式
func Trace(v ...interface{}) {
	if trace == nil {
		return
	}

	trace.Println(v...)
}

// Tracef 相当于 TRACE().Printf(format, v...) 的简写方式
func Tracef(format string, v ...interface{}) {
	if trace == nil {
		return
	}

	trace.Printf(format, v...)
}

// 获取 WARN 级别的 log.Logger 实例，在未指定 warn 级别的日志时，该实例返回一个 nil。
func WARN() *log.Logger {
	return warn
}

// Warn 相当于 WARN().Println(v...) 的简写方式
func Warn(v ...interface{}) {
	if warn == nil {
		return
	}

	warn.Println(v...)
}

// Warnf 相当于 WARN().Printf(format, v...) 的简写方式
func Warnf(format string, v ...interface{}) {
	if warn == nil {
		return
	}

	warn.Printf(format, v...)
}

// 获取 ERROR 级别的 log.Logger 实例，在未指定 error 级别的日志时，该实例返回一个 nil。
func ERROR() *log.Logger {
	return erro
}

// Error 相当于 ERROR().Println(v...) 的简写方式
func Error(v ...interface{}) {
	if erro == nil {
		return
	}

	erro.Println(v...)
}

// Errorf 相当于 ERROR().Printf(format, v...) 的简写方式
func Errorf(format string, v ...interface{}) {
	if erro == nil {
		return
	}

	erro.Printf(format, v...)
}

// 获取 CRITICAL 级别的 log.Logger 实例，在未指定 critical 级别的日志时，该实例返回一个 nil。
func CRITICAL() *log.Logger {
	return critical
}

// Critical 相当于 CRITICAL().Println(v...)的简写方式
func Critical(v ...interface{}) {
	if critical == nil {
		return
	}

	critical.Println(v...)
}

// Criticalf 相当于 CRITICAL().Printf(format, v...) 的简写方式
func Criticalf(format string, v ...interface{}) {
	if critical == nil {
		return
	}

	critical.Printf(format, v...)
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

// 输出错误信息，然后退出程序。
func Fatal(v ...interface{}) {
	All(v...)
	Flush()
	os.Exit(2)
}

// 输出错误信息，然后退出程序。
func Fatalf(format string, v ...interface{}) {
	Allf(format, v...)
	Flush()
	os.Exit(2)
}

// 输出错误信息，然后触发 panic。
func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	All(s)
	Flush()
	panic(s)
}

// 输出错误信息，然后触发 panic。
func Panicf(format string, v ...interface{}) {
	Allf(format, v...)
	Flush()
	panic(fmt.Sprintf(format, v...))
}
