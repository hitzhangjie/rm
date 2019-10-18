package main

import (
	"fmt"
	"github.com/hitzhangjie/rm/cmd"
	"os"
)

func main() {

	helpCmd := cmd.SubCmd("help")

	args := os.Args[1:]
	if len(args) == 0 {
		if err := helpCmd.Run(nil); err != nil {
			fmt.Printf("subcmd `help` error: %v\n", err)
		}
		return
	}

	subcmd := cmd.SubCmd(os.Args[1])
	if subcmd != nil {
		if err := subcmd.Run(os.Args[2:]); err != nil {
			fmt.Printf("subcmd `%s` error: %v\n", os.Args[1], err)
		}
		return
	}

	rmCmd := cmd.NewRmCmd()
	if err := rmCmd.Run(os.Args[1:]); err != nil {
		fmt.Printf("rm error: %v\n", err)
	}
}
