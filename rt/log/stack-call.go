package log

import (
	"fmt"
	"runtime"
	"strings"
)

// CallStack 调用栈
type CallStack struct {
	skipPackages []string
}

type errStack struct {
	call    *CallStack
	stackPC []uintptr
	raw     error
}

func (e *errStack) Error() string {
	return e.call.error(e)
}

// NewCallStack 创建调用栈
func NewCallStack() *CallStack {
	return &CallStack{
		skipPackages: []string{
			"runtime.",
			// "runtime.main",
			// "runtime.goexit",
			// "runtime.doInit",
			// "github.com/sirupsen/logrus",
			// "github.com/seerx/base.ErrorStack",
			// "github.com/seerx/base/pkg/logs",
			// "github.com/seerx/base.(*CallStack).WrapError",
		},
	}
}

// AddSkipPackage 添加跳过包
func (c *CallStack) AddSkipPackage(pkg string) {
	c.skipPackages = append(c.skipPackages, pkg)
}

// WrapError 错误堆栈信息
func (c *CallStack) WrapError(err error) error {
	// pcs := make([]uintptr, 32)
	// // skip func StackError invocations
	// count := runtime.Callers(2, pcs)
	// return &errStack{
	// 	raw:     err,
	// 	stackPC: pcs[:count],
	// }
	return c.WrapErrorSkip(err, 0)
}

// WrapErrorSkip 错误堆栈信息
func (c *CallStack) WrapErrorSkip(err error, skip int) error {
	pcs := make([]uintptr, 32)
	// skip func StackError invocations
	count := runtime.Callers(2+skip, pcs)
	return &errStack{
		call:    c,
		raw:     err,
		stackPC: pcs[:count],
	}
}

func (c *CallStack) containsPackage(Function string) bool {
	for _, pkg := range c.skipPackages {
		if strings.Index(Function, pkg) == 0 {
			return true
		}
	}
	return false
}

func (c *CallStack) error(e *errStack) string {
	frames := runtime.CallersFrames(e.stackPC)

	var (
		f     runtime.Frame
		more  bool
		index int
	)

	errStr := strings.Builder{}
	// errString := ""
	if e.raw != nil && e.raw.Error() != "" {
		errStr.WriteString(e.raw.Error() + "\n")
		// errString = e.raw.Error() + "\n"
	}

	for {
		f, more = frames.Next()
		if index = strings.Index(f.File, "src"); index != -1 {
			// trim GOPATH or GOROOT prifix
			f.File = string(f.File[index+4:])
		}
		if !c.containsPackage(f.Function) { // 不要跳过的包
			frame := fmt.Sprintf("  %s\n    %s:%d\n", f.Function, f.File, f.Line)
			errStr.WriteString(frame)
			// errString = fmt.Sprintf("%s%s\n\t%s:%d\n", errString, f.Function, f.File, f.Line)
		}
		if !more {
			break
		}
	}
	return errStr.String()
}
