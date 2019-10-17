package cmd

import (
	"flag"
)

var cmds = []Commander{
	NewHelpCmd(),
	NewPinCmd(),
	newUnpinCmd(),
}

type Commander interface {
	Name() string
	Description() string
	FlagSet() *flag.FlagSet
	Run(args []string) error
}

type baseCmd struct {
	description string
	flagSet     *flag.FlagSet
}

func (c *baseCmd) Description() string {
	return c.description
}

func (c *baseCmd) FlagSet() *flag.FlagSet {
	return c.flagSet
}

func SubCmd(name string) Commander {
	for _, v := range cmds {
		if v.Name() == name {
			return v
		}
	}
	return nil
}
