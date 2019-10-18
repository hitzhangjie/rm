package cmd

import (
	"flag"
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
		description: `
rm [options] <target>...:
	same as shell /bin/rm`,
		flagSet:     newRmCmdFlagSet(),
	}
	return &RmCmd{c}
}

func newRmCmdFlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("rm", flag.PanicOnError)

	fs.Bool("f", false, "ignore nonexistent files and arguments, never prompt")
	fs.Bool("force", false, "ignore nonexistent files and arguments, never prompt")

	fs.Bool("i", false, "prompt before every removal")
	fs.Bool("I", false, `
prompt once before removing more than three files, or when  removing  recur-
sively;  less  intrusive than -i, while still giving protection against most
mistakes`)

	fs.String("interactive", "WHEN", `
prompt according to WHEN: never, once (-I), or always  (-i);  without  WHEN,
prompt always`)
	fs.Bool("one-file-system", false, `
when  removing a hierarchy recursively, skip any directory that is on a file
system different from that of the corresponding command line argument`)

	fs.Bool("no-preserve-root", false, "do not treat '/' specially")
	fs.Bool("preserve-root", false, `
do not remove '/' (default); with 'all', reject any command line argument on
a separate device from its parent`)

	fs.Bool("r", false, "remove directories and their contents recursively")
	fs.Bool("R", false, "remove directories and their contents recursively")

	fs.Bool("d", false, "remove empty directories")
	fs.Bool("v", false, "explain what is being done")

	fs.String("help", "", "display this help and exit")
	fs.String("version", "", "output version information and exit")

	return fs
}

func (c *RmCmd) Run(args []string) error {

	c.flagSet.Parse(args)

	targets := c.flagSet.Args()

	for _, target := range targets {

		fmt.Println(target)

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
	}
	return nil
}
