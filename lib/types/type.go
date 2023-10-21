package types

import (
	"bytes"
	"fmt"

	"github.com/lspaccatrosi16/go-cli-tools/input"
)

type cmd struct {
	Name        string
	Description string
	Exec        *func() error
}

func (c *cmd) Run() error {
	if c == nil {
		panic("exec property of command must not be nil")
	}
	return (*c.Exec)()
}

type manager struct {
	cmds []*cmd
}

func (m *manager) Help() {
	maxCmdLength := 0
	cmds := []string{}
	descriptions := []string{}
	for _, cmd := range m.cmds {
		if len(cmd.Name) > maxCmdLength {
			maxCmdLength = len(cmd.Name)
		}
		cmds = append(cmds, cmd.Name)
		descriptions = append(descriptions, cmd.Description)
	}
	buf := bytes.NewBuffer(nil)

	for i := 0; i < len(cmds); i++ {
		cmd := cmds[i]
		desc := descriptions[i]

		fmt.Fprintf(buf, "%-*s: %s\n", maxCmdLength+2, cmd, desc)
	}

	fmt.Println(buf.String())
}

func (m *manager) Register(name string, description string, exec func() error) {
	newcmd := cmd{
		Name:        name,
		Description: description,
		Exec:        &exec,
	}
	m.cmds = append(m.cmds, &newcmd)
}

func (m *manager) Run(str string) {
	for _, cmd := range m.cmds {
		if cmd.Name == str {
			err := cmd.Run()
			if err != nil {
				fmt.Println("an error was encountered whilst running the command:\n", err.Error())
			}
			return
		}
	}

	fmt.Printf("command \"%s\" was not found\n", str)
	m.Help()
}

func (m *manager) Gui() {
	options := []input.SelectOption{}
	for _, cmd := range m.cmds {
		options = append(options, input.SelectOption{Name: fmt.Sprintf("%s: %s", cmd.Name, cmd.Description), Value: cmd.Name})
	}
	selected, err := input.GetSearchableSelection("Select the command to execute", options)
	if err != nil {
		panic(err)
	}
	m.Run(selected)
}

func NewManager() manager {
	return manager{}
}
