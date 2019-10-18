package cmd

import (
	//"flag"
	flag "github.com/spf13/pflag"
	"os"
	"path/filepath"
)

type UnpinCmd struct {
	baseCmd
}

func newUnpinCmd() *UnpinCmd {

	cmd := baseCmd{
		description: `
rm unpin: 
	-r, unpin target recursively`,
		flagSet:     newUnpinCmdFlagSet(),
	}

	return &UnpinCmd{cmd}
}

func newUnpinCmdFlagSet() *flag.FlagSet {

	fs := flag.NewFlagSet("unpin", flag.PanicOnError)
	fs.Bool("r", false, "recursively if target is directory")

	return fs
}

func (c *UnpinCmd) Name() string {
	return "unpin"
}

// Run 执行pin动作
//
// rm pin
// rm pin -r
// rm pin -r .
// rm pin -r target
// rm pin target
func (c *UnpinCmd) Run(args []string) error {

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

	target := c.flagSet.Args()

	// cwd
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if len(target) == 0 {
		target = []string{cwd}
	}

	for _, t := range target {
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
		// check -r when unpin directory, while unpin single file, we don't need to check -r
		if !fin.IsDir() {
			dir := filepath.Dir(abs)
			err := os.RemoveAll(filepath.Join(dir, pinLock))
			if err != nil {
				return err
			}
			return nil
		}

		if !recursive {
			err := os.RemoveAll(filepath.Join(abs, pinLock))
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
			err = os.RemoveAll(filepath.Join(entry, pinLock))
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
