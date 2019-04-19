//
// Date: 2018-03-20
// Author: spicer (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-12-28
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package services

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mgutz/ansi"
)

//
// Info Log. Used to log inportant information but a human does not need to review unless they are debugging something.
//
func Info(err error) {
	log.Println("[App:Info] " + MyCaller() + " : " + ansi.Color(err.Error(), "magenta"))
}

//
// Allow us to pass in just text instead of an error. Not always do we have a err to pass in.
//
func InfoMsg(msg string) {
	Info(errors.New(msg))
}

//
// Critical - We used this when we want to make splash. All Critical errors should be reviewed by a human.
//
func Critical(err error) {
	log.Println(ansi.Color("[App:Critical] "+MyCaller()+" : "+err.Error(), "yellow"))
}

//
// Fatal Log. We use this wehn the app should die and not continue running.
//
func Fatal(err error) {
	log.Fatal(ansi.Color("[App:Fatal] "+MyCaller()+" : "+err.Error(), "red"))
}

// ----- TODO(spicer): Get rid of functions below --- //

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
	parts := strings.Split(file, "app.skyclerk.com")

	if len(parts) == 2 {
		filePath = "app.skyclerk.com" + parts[1]
	} else {
		filePath = filepath.Base(file)
	}

	return fmt.Sprintf("%s:%d %s", filePath, line, fnName)
}

/* End File */
