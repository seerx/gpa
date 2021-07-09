package log

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/seerx/gpa/rt/log/transfers"
	"github.com/sirupsen/logrus"
)

const logruspkg = "github.com/sirupsen/logrus"

// Builder 日志 Builder
type Builder struct {
	appTag          string       // 应用标志
	prettyJSON      bool         // 输出格式化的 json
	timestampFormat string       // 日期格式
	reportCaller    bool         // 是否输出日志发生地址 , 文件 函数 行号
	level           logrus.Level // 日志输出级别
	outputJSON      bool

	console      bool     // 在控制台输出
	udpHost      string   // 接收日志的 udp 主机
	udpPort      int      // 接收日志的 udp 端口
	skipPackages []string // 忽略的包
}

// NewBuilder 创建 builder
func NewBuilder() *Builder {
	return &Builder{
		console:         true,
		level:           logrus.InfoLevel,
		timestampFormat: time.RFC3339,
		reportCaller:    false,
	}
}

// Build 创建日志
func (b *Builder) Build() *log.Logger {
	// 创建堆栈调用生成组件
	added := false
	stack := NewCallStack()
	for _, pkg := range b.skipPackages {
		stack.AddSkipPackage(pkg)
		if pkg == logruspkg && !added {
			added = true
		}
	}
	if !added {
		stack.AddSkipPackage(logruspkg)
	}

	// 创建日志
	var logger = log.New(nil, "", 0)
	setNull(logger)
	// if b.outputJSON {
	// 	logger.SetOutput()
	// 	logger.Formatter = &logrus.JSONFormatter{
	// 		TimestampFormat:  b.timestampFormat,
	// 		DisableTimestamp: false,
	// 		DataKey:          "",
	// 		FieldMap:         nil,
	// 		CallerPrettyfier: nil,
	// 		PrettyPrint:      b.prettyJSON,
	// 	}
	// } else {
	// 	logger.Formatter = &TextFormatter{
	// 		timeFormat: b.timestampFormat,
	// 	}
	// }
	logger.Level = b.level
	// logger.ReportCaller = b.reportCaller

	var txfns []transfers.TransferFn
	if b.console {
		txfns = append(txfns, MakeTransfer(nil))
	}
	if b.udpHost != "" && b.udpPort > 0 {
		txfns = append(txfns, MakeTransfer(&transfers.TransferConfigure{
			Type:   transfers.UDP,
			Server: b.udpHost,
			Port:   b.udpPort,
		}))
	}
	// logger.
	logger.AddHook(NewTransferHook(b.appTag,
		stack,
		logger.Formatter,
		txfns...))

	return logger
}

// ReportCaller 是否报告日志地址
func (b *Builder) ReportCaller(report bool) *Builder {
	b.reportCaller = report
	return b
}

// Level 日志级别
func (b *Builder) Level(level logrus.Level) *Builder {
	b.level = level
	return b
}

// WriteToUDP 设置日志输出到 udp
func (b *Builder) WriteToUDP(host string, port int) *Builder {
	b.udpHost = host
	b.udpPort = port
	return b
}

// WriteToConsole 是否输出到控制台
func (b *Builder) WriteToConsole(write bool) *Builder {
	b.console = write
	return b
}

// AddSkipPackages 添加错误堆栈中忽略的包
func (b *Builder) AddSkipPackages(pkgs ...string) *Builder {
	b.skipPackages = append(b.skipPackages, pkgs...)
	return b
}

// OutputJSON 输出 json 格式
func (b *Builder) OutputJSON(json bool) *Builder {
	b.outputJSON = json
	return b
}

// PrettyJSON 是否格式化输出
func (b *Builder) PrettyJSON(pretty bool) *Builder {
	b.prettyJSON = pretty
	return b
}

// TimestampFormat 设置时间格式
func (b *Builder) TimestampFormat(format string) *Builder {
	b.timestampFormat = format
	return b
}

// AppTag 设置应用标志
func (b *Builder) AppTag(appTag string) *Builder {
	b.appTag = appTag
	return b
}

func setNull(logger *log.Logger) {
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	writer := bufio.NewWriter(src)
	logger.SetOutput(writer)
}

// MakeTransfer 创建转发函数
func MakeTransfer(cfg *transfers.TransferConfigure) transfers.TransferFn {
	if cfg == nil || cfg.Type == transfers.CONSOLE {
		return transfers.CreateConsoleTransfer(cfg)
	}
	if cfg.Type == transfers.UDP {
		return transfers.CreateUDPTransfer(cfg)
	}
	panic("未实现的转发")
}
