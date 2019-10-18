package cmd

import (
	"fmt"
	//"flag"
	flag "github.com/spf13/pflag"
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
		flagSet: newRmCmdFlagSet(),
	}
	return &RmCmd{c}
}

func newRmCmdFlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("rm", flag.ContinueOnError)

	fs.BoolP("force", "f", false, "ignore nonexistent files and arguments, never prompt")

	fs.BoolP("i", "i", false, "prompt before every removal")
	fs.BoolP("I", "I", false, `
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

	fs.BoolP("r", "r", false, "remove directories and their contents recursively")
	fs.BoolP("R", "R", false, "remove directories and their contents recursively")

	fs.BoolP("d", "d", false, "remove empty directories")
	fs.BoolP("v", "v", false, "explain what is being done")

	fs.String("help", "", "display this help and exit")
	fs.String("version", "", "output version information and exit")

	return fs
}

func (c *RmCmd) Run(args []string) error {

	err := c.flagSet.Parse(args)
	if err != nil {
		return err
	}

	targets := c.flagSet.Args()

	for _, target := range targets {
		//println(target)
		abs := target

		// convert target to absolute path
		if !filepath.IsAbs(target) {
			if dir, err := os.Getwd(); err != nil {
				return err
			} else {
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
			// firstly, check parent directory
			pdir := filepath.Dir(abs)
			p := filepath.Join(pdir, pinLock)
			if fin, err := os.Lstat(p); err == nil && !fin.IsDir() {
				//return fmt.Errorf("directory %s and files underneath, pined by %s", cwd, p)
				return fmt.Errorf("pined: directory %s and files underneath", pdir)
			}

			// secondly, check target directory
			err := filepath.Walk(abs, func(entry string, fin os.FileInfo, err error) error {
				if strings.Contains(fin.Name(), pinLock) {
					return fmt.Errorf("pined: directory %s and files underneath", filepath.Dir(entry))
				}
				return nil
			})
			if err != nil {
				return err
			}
		} else {
			dir := filepath.Dir(abs)
			p := filepath.Join(dir, pinLock)
			if fin, err := os.Lstat(p); err == nil && !fin.IsDir() {
				//return fmt.Errorf("directory %s and files underneath, pined by %s", cwd, p)
				return fmt.Errorf("pined: directory %s and files underneath", dir)
			}
		}

		// call /bin/rm to delete
		//
		//cmd := exec.Command("/bin/rm", args...)
		//if output, err := cmd.CombinedOutput(); err != nil {
		//	return fmt.Errorf("%s, %s", err, string(output))
		//} else {
		//	fmt.Println(string(output))
		//}
	}

	// all check passed, then call /bin/rm to delete
	// this behavior is a little different from /bin/rm, but it's much secure
	cmd := exec.Command("/bin/rm", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s, %s", err, string(output))
	} else {
		fmt.Println(string(output))
	}

	return nil
}
