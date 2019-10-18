package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type RmCmd struct {
	baseCmd
}

func NewRmCmd() *RmCmd {
	c := baseCmd{
		description: "rm -r -f -i -... target",
		flagSet:     nil,
	}
	return &RmCmd{c}
}

func (c *RmCmd) Run(args []string) error {
	n := len(args)
	target := args[n-1]

	var (
		cwd string
		abs = target
	)

	// convert target to absolute path
	if !filepath.IsAbs(target) {
		if dir, err := os.Getwd(); err != nil {
			return err
		} else {
			cwd = dir
			abs = filepath.Join(dir, target)
		}
	}

	// if target not exist, fail fast
	fin, err := os.Lstat(abs)
	if err != nil {
		return err
	}

	// if target is dir, check whether `pinlock` exists underneath recursively
	// if target is file, check whether `pinlock` dir(target) exists
	//
	// if one `pinlock` found, that means some files is protected, fail
	// if no `pinlock`s found, call /bin/rm to finish the task
	if fin.IsDir() {
		err := filepath.Walk(cwd, func(entry string, fin os.FileInfo, err error) error {
			if strings.Contains(fin.Name(), pinLock) {
				return fmt.Errorf("pined: directory %s and files underneath", filepath.Dir(entry))
			}
			return nil
		})
		if err != nil {
			return err
		}
	} else {
		p := filepath.Join(cwd, pinLock)
		if fin, err := os.Lstat(p); err == nil && !fin.IsDir() {
			//return fmt.Errorf("directory %s and files underneath, pined by %s", cwd, p)
			return fmt.Errorf("pined: directory %s and files underneath", cwd)
		}
	}

	// call /bin/rm to delete
	cmd := exec.Command("/bin/rm", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s, %s", err, string(output))
	} else {
		fmt.Println(string(output))
	}
	return nil
}
