//
// Date: 2018-03-20
// Author: spicer (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-12-28
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package services

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mgutz/ansi"
)

//
// Normal Log.
//
func LogInfo(message string) {
	log.Println(message)
}

//
// Debug Log.
//
func LogDebug(message string) {
	log.Println(message)
}

//
// Fatal Log.
//
func LogFatal(err error) {
	log.Fatal(err)
}

//
// Warning Log.
//
func LogWarning(err error) {

	caller := MyCaller()

	// Standard out
	log.Println(ansi.Color("[App:Warning] "+caller+" : "+err.Error(), "yellow"))
}

//
// Error Log.
//
func Error(err error) {

	caller := MyCaller()

	// Standard out
	log.Println(ansi.Color("[App:Error] "+caller+" : "+err.Error(), "red"))
}

//
// MyCaller returns the caller of the function that called the logger :)
//
func MyCaller() string {
	var filePath string
	var fnName string

	pc, file, line, ok := runtime.Caller(3)

	if !ok {
		file = "?"
		line = 0
	}

	fn := runtime.FuncForPC(pc)

	if fn == nil {
		fnName = "?()"
	} else {
		dotName := filepath.Ext(fn.Name())
		fnName = strings.TrimLeft(dotName, ".") + "()"
	}

	// Make the base of this code.
	parts := strings.Split(file, "app.options.cafe")

	if len(parts) == 2 {
		filePath = "app.options.cafe" + parts[1]
	} else {
		filePath = filepath.Base(file)
	}

	return fmt.Sprintf("%s:%d %s", filePath, line, fnName)
}

/* End File */
