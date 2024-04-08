package xflag

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var defaultCommandSet = NewCommandSet()

var (
	commandRunName  = "run"
	commandExecName = "exec"
	commandHelpName = "help"
)

func NewCommandSet() *CommandSet {
	m := &CommandSet{
		commands: make(map[string]*flag.FlagSet),
	}
	m.Init()
	return m
}

type CommandSet struct {
	commands     map[string]*flag.FlagSet
	parsed       bool
	args         []string
	commandName  string
	flagAppName  string
	flagHttpPort string
	output       io.Writer
	terminal     bool

	Help func() error
}

func (m *CommandSet) Init() {

	m.CommandLine(commandRunName).StringVar(&m.flagAppName, "app-name", "", "")
	m.CommandLine(commandRunName).StringVar(&m.flagHttpPort, "http-port", "", "")
	m.CommandLine(commandExecName).StringVar(&m.flagAppName, "app-name", "", "")
	m.CommandLine(commandExecName).StringVar(&m.flagHttpPort, "http-port", "", "")
	m.CommandLine(commandHelpName)

	m.Help = func() error { return printHelp(m.Output(), false) }
}

func (m *CommandSet) Args() []string {
	return m.args
}

func (m *CommandSet) Parse(args []string) bool {
	if m.parsed {
		return m.parsed
	}
	m.args = args
	if len(args) < 2 {
		return m.parsed
	}
	fs, ok := m.commands[args[1]]
	if !ok {
		return m.parsed
	}
	if err := fs.Parse(args[2:]); err != nil {
		return m.parsed
	}
	m.parsed = true
	m.commandName = args[1]
	return m.parsed
}

func (m *CommandSet) Parsed() bool {
	return m.parsed
}

func (m *CommandSet) SetOutput(output io.Writer) {
	m.output = output
}

func (m *CommandSet) Output() io.Writer {
	if m.output == nil {
		return os.Stderr
	}
	return m.output
}

func (m *CommandSet) CommandName() string {
	return m.commandName
}

func (m *CommandSet) CommandLine(command string) *flag.FlagSet {
	fs, ok := m.commands[command]
	if !ok {
		fs = flag.NewFlagSet(command, flag.ExitOnError)
		fs.SetOutput(m.Output())
		m.commands[command] = fs
	}
	return fs
}

func (m *CommandSet) Service(command string, f func() error) error {
	if m.terminal {
		return nil
	}
	if !m.Parse(m.Args()) {
		printHelp(m.Output(), true)
		return nil
	}
	if f == nil {
		return nil
	}

	switch m.CommandName() {
	case "help":
		m.terminal = true
		return m.Help()
	}

	if !strings.EqualFold(command, m.CommandName()) {
		return nil
	}
	return f()
}

func (m *CommandSet) AppName() string {
	return m.flagAppName
}

func (m *CommandSet) HttpPort() string {
	return m.flagHttpPort
}

func printHelp(output io.Writer, exitOnError bool) error {
	s := `
	app run
	app exec
	app help`
	fmt.Fprintln(output, s)
	if exitOnError {
		os.Exit(2)
	}
	return nil
}

func Parse() bool {
	return defaultCommandSet.Parse(os.Args)
}

func Parsed() bool {
	return defaultCommandSet.Parsed()
}

func CommandLine(command string) *flag.FlagSet {
	return defaultCommandSet.CommandLine(command)
}

func Service(command string, f func() error) error {
	Parse()
	return defaultCommandSet.Service(command, f)
}

func SetOutput(output io.Writer) {
	defaultCommandSet.SetOutput(output)
}

func Output() io.Writer {
	return defaultCommandSet.Output()
}

func AppName() string {
	return defaultCommandSet.AppName()
}

func HttpPort() string {
	return defaultCommandSet.HttpPort()
}

func Help() error {
	return defaultCommandSet.Help()
}
