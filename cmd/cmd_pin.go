package cmd

import (
	//"flag"
	flag "github.com/spf13/pflag"
	"os"
	"path/filepath"
)

const pinLock = ".pinLock"

type PinCmd struct {
	baseCmd
}

func NewPinCmd() *PinCmd {

	cmd := baseCmd{
		description: `
rm pin: 
	-r, pin target recursively`,
		flagSet: newPinCmdFlagSet(),
	}

	return &PinCmd{cmd}
}

func newPinCmdFlagSet() *flag.FlagSet {

	fs := flag.NewFlagSet("pin", flag.PanicOnError)
	fs.BoolP("r", "r", false, "recursively if target is directory")

	return fs
}

func (c *PinCmd) Name() string {
	return "pin"
}

// Run 执行pin动作
//
// rm pin
// rm pin -r
// rm pin -r .
// rm pin -r target
// rm pin target
func (c *PinCmd) Run(args []string) error {

	if err := c.flagSet.Parse(args); err != nil {
		return err
	}

	//println(c.flagSet.NFlag()) // 命令行解析时传递的命令行选项正确设置的个数（默认值的不算）
	//println(c.flagSet.NArg()) // 命令行定义的参数解析完成后还剩下什么参数，比如rm pin hello，此时为1（hello），比如rm pin hello world，此时为2

	//recursive := c.flagSet.Lookup("r").Value.(flag.Getter).Get().(bool)
	recursive, err := c.flagSet.GetBool("r")
	if err != nil {
		return err
	}

	// cwd
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	targets := c.flagSet.Args()

	for _, t := range targets {
		// check whether `t` is absolute path
		abs := t
		if !filepath.IsAbs(t) {
			abs = filepath.Join(cwd, t)
		}
		// check whether `t` is directory or not
		fin, err := os.Lstat(abs)
		if err != nil {
			return err
		}
		// check -r when pin directory
		if !fin.IsDir() {
			dir := filepath.Dir(abs)
			_, err := os.Create(filepath.Join(dir, pinLock))
			if err != nil {
				return err
			}
			return nil
		}

		if !recursive {
			_, err := os.Create(filepath.Join(abs, pinLock))
			if err != nil {
				return err
			}
			return nil
		}

		err = filepath.Walk(abs, func(entry string, fin os.FileInfo, err error) error {
			if os.IsNotExist(err) {
				return nil
			}
			if err != nil {
				return err
			}
			if !fin.IsDir() {
				return nil
			}
			_, err = os.Create(filepath.Join(entry, pinLock))
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}
