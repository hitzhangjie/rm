package cmd

import (
	"fmt"
	"strings"
)

type HelpCmd struct {
	baseCmd
}

func NewHelpCmd() *HelpCmd {
	c := baseCmd{
		description: `
rm help: 
	display help info`,
		flagSet: nil,
	}
	return &HelpCmd{c}
}

func (c *HelpCmd) Name() string {
	return "help"
}

func (c *HelpCmd) Run(args []string) error {
	buf := strings.Builder{}

	for _, v := range cmds {
		buf.WriteString(v.Description() + "\n")
	}
	fmt.Println(buf.String())
	return nil
}
