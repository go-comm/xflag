package xflag

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var Commands = NewXFlag(os.Args[0])

func Parse() {
	Commands.Parse(os.Args[1:])
}

func CommandLine(command string) *flag.FlagSet {
	return Commands.CommandLine(command)
}

func CommandLineRun() *flag.FlagSet {
	return Commands.CommandLine("run")
}

func CommandLineExec() *flag.FlagSet {
	return Commands.CommandLine("exec")
}

func Service(command string, fn func() error) error {
	return Commands.Service(command, fn)
}

func Run(command string, fn func() error) error {
	return Commands.Service("run", fn)
}

func Exec(command string, fn func() error) error {
	return Commands.Service("exec", fn)
}

var Usage = func() {
	fmt.Fprintf(Commands.Output(), "Usage of %s:\n", os.Args[0])
	Commands.Print()
}

func NewXFlag(name string) *XFlag {
	xf := &XFlag{
		name: name,
	}
	return xf
}

type XFlag struct {
	name         string
	commandLines map[string]*flag.FlagSet
	args         []string
	parsed       bool

	Usage  func()
	output io.Writer
}

func (xf *XFlag) Output() io.Writer {
	if xf.output == nil {
		return os.Stderr
	}
	return xf.output
}

func (xf *XFlag) CommandLine(command string) *flag.FlagSet {
	if xf.commandLines == nil {
		xf.commandLines = make(map[string]*flag.FlagSet)
	}
	command = fmtCommand(command)
	fs, ok := xf.commandLines[command]
	if !ok {
		fs = flag.NewFlagSet(command, flag.ExitOnError)
		fs.SetOutput(xf.Output())
		xf.commandLines[command] = fs
	}
	return fs
}

func (xf *XFlag) Parse(args []string) error {
	xf.parsed = true
	if len(args) <= 1 {
		xf.Usage()
		return errors.New("xflag: args length got <1, want >1")
	}
	xf.args = args
	fs, ok := xf.commandLines[args[0]]
	if !ok {
		xf.Usage()
		return nil
	}
	err := fs.Parse(args[1:])
	if err != nil {
		return xf.failError(args[1], err)
	}
	return nil
}

func (xf *XFlag) Service(command string, fn func() error) error {
	if fn == nil {
		return nil
	}
	command = fmtCommand(command)
	if _, ok := xf.commandLines[command]; !ok {
		xf.Usage()
		return xf.failError(command, errors.New("no command"))
	}
	return fn()
}

func (xf *XFlag) failError(command string, err error) error {
	return fmt.Errorf("xflag: command '%s', %w", command, err)
}

func (xf *XFlag) Print() {
	for _, v := range xf.commandLines {
		v.PrintDefaults()
	}
}

func fmtCommand(command string) string {
	return strings.ToLower(command)
}
