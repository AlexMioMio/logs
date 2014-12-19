// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/issue9/logs"
)

func main() {
	err := logs.InitFromFile("./config.xml")
	if err != nil {
		//panic(err)
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	defer logs.Flush()

	logs.Info("INFO1")
	logs.Debugf("DEBUG %v", 1)
}
