// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"bytes"
	"io"
	"log"
	"testing"

	"github.com/issue9/assert"
)

var (
	debugW    = new(bytes.Buffer)
	infoW     = new(bytes.Buffer)
	errorW    = new(bytes.Buffer)
	traceW    = new(bytes.Buffer)
	warnW     = new(bytes.Buffer)
	criticalW = new(bytes.Buffer)
)

func resetLog(t *testing.T) {
	a := assert.New(t)

	infoW.Reset()
	debugW.Reset()
	errorW.Reset()
	traceW.Reset()
	warnW.Reset()
	criticalW.Reset()

	a.True(infoW.Len() == 0)
	a.True(debugW.Len() == 0)
	a.True(errorW.Len() == 0)
	a.True(traceW.Len() == 0)
	a.True(warnW.Len() == 0)
	a.True(criticalW.Len() == 0)

	info = log.New(infoW, "[INFO]", log.LstdFlags)
	debug = log.New(debugW, "[DEBUG]", log.LstdFlags)
	erro = log.New(errorW, "[ERROR]", log.LstdFlags)
	trace = log.New(traceW, "[TRACE]", log.LstdFlags)
	warn = log.New(warnW, "[WARN]", log.LstdFlags)
	critical = log.New(criticalW, "[CRITICAL]", log.LstdFlags)
}

func checkLog(t *testing.T) {
	a := assert.New(t)

	a.True(infoW.Len() > 0)
	a.True(debugW.Len() > 0)
	a.True(errorW.Len() > 0)
	a.True(traceW.Len() > 0)
	a.True(warnW.Len() > 0)
	a.True(criticalW.Len() > 0)
}

func TestAll(t *testing.T) {
	resetLog(t)
	All("abc")
	checkLog(t)
}

func TestAllf(t *testing.T) {
	resetLog(t)
	Allf("abc")
	checkLog(t)
}

func debugWInit(args map[string]string) (io.Writer, error) {
	return debugW, nil
}

func TestInitFormXMLString(t *testing.T) {
	a := assert.New(t)

	// 重新注册以下用到的writer
	clearInitializer()
	a.True(Register("debug", logContInitializer), "注册debug时失败")
	a.True(Register("buffer", bufferInitializer), "注册buffer时失败")
	a.True(Register("debugW", debugWInit), "注册debugW时失败")

	xml := `
<?xml version="1.0" encoding="utf-8" ?>
<logs>
	<debug prefix="[DEBUG]">
		<buffer size="10">
			<debugW />
		</buffer>
	</debug>
</logs>
`
	debugW.Reset()
	conts.Add(infoW) // 触发initFromXmlString中的重置功能
	a.True(conts.Len() == 1)
	a.NotError(InitFromXMLString(xml))
	a.True(critical == nil)   // InitFromXMLString会重置所有的日志指向
	a.True(CRITICAL() == nil) // InitFromXMLString会重置所有的日志指向

	Debug("abc")
	a.True(debugW.Len() == 0) // 缓存未达10，依然为空
	Allf("def\n")
	a.True(debugW.Len() == 0) // 缓存未达10，依然为空

	// 测试Flush
	Flush()
	a.True(debugW.Len() > 0)
}
