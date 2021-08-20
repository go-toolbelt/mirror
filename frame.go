package mirror

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"
)

type Frame struct {
	Package   string
	Function  string
	File      string
	Line      int
	Formatted string
}

func fromRuntimeFrame(frame runtime.Frame) Frame {
	packageName, functionName := splitFunction(frame.Function)
	fileName := toFileName(frame.File)
	line := frame.Line
	return Frame{
		Package:   packageName,
		Function:  functionName,
		File:      fileName,
		Line:      line,
		Formatted: toFormatted(packageName, functionName, fileName, line),
	}
}

var re = regexp.MustCompile(`[\*\(\)]`)

func splitFunction(function string) (packageName string, functionName string) {
	i := strings.LastIndex(function, "/") + 1
	j := strings.Index(function[i:], ".")
	return function[:i+j], re.ReplaceAllString(function[i+j+1:], "")
}

func toFileName(file string) string {
	return file[strings.LastIndex(file, "/")+1:]
}

func toFormatted(
	packageName string,
	functionName string,
	fileName string,
	line int,
) string {
	return fmt.Sprintf("\tat %s.%s(%s:%d)", packageName, functionName, fileName, line)
}

func (frame Frame) String() string {
	return frame.Formatted
}
