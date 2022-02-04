package utils

import (
	"fmt"
	"runtime"
)

func Fname() string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)

	return runtime.FuncForPC(pc[0]).Name()
}

func PrintCallerPos() {
	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	fmt.Printf("%s:%d %s\n", file, line, f.Name())
}

func GetCallerPos() string {
	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	return fmt.Sprintf("%s:%d %s\n", file, line, f.Name())
}
