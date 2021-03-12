package commons

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
)

const (
	ERROR_SPLITTER = " | "
)

// TODO: don't use fmt? move this to utils
func WrapError(innerErr error, outerErrString string) error {
	outerErrString = "%w" + ERROR_SPLITTER + outerErrString
	return fmt.Errorf(outerErrString, innerErr)
}

func GetCallerDetails() string {
	programCounter := make([]uintptr, 1)
	// Skipping two frames
	n := runtime.Callers(3, programCounter)
	frames := runtime.CallersFrames(programCounter[:n])
	frame, _ := frames.Next()
	return frame.Function + ":" + strconv.Itoa(frame.Line)
}

func NameOfFunction(function interface{}) string {
	value := reflect.ValueOf(function)
	if value.Kind() == reflect.Func {
		if rf := runtime.FuncForPC(value.Pointer()); rf != nil {
			return rf.Name()
		}
	}
	return value.String()
}
