package main

import (
    "fmt"
    "os"
	"io"
    "flag"
    "sync"
    "text/template"
    "github.com/golang/glog"
    "github.com/lshengjian/kms/util"


)
var exitStatus = 0
var exitMu sync.Mutex
var commands = []*Command{
    cmdWeb,
	cmdDownload,
	cmdVersion,
}



func setExitStatus(n int) {
	exitMu.Lock()
	if exitStatus < n {
		exitStatus = n
	}
	exitMu.Unlock()
}

var usageTemplate = `
KMS: store billions of files and serve them fast!

Usage:

	kms command [arguments]

The commands are:
{{range .}}
    {{.Name | printf "%-11s"}} {{.Short}}{{end}}

Use "kms help [command]" for more information about a command.

`

func main() {
	flag.Usage = usage
    flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "logs")
	flag.Set("v", "4")
	flag.Parse()
    args := flag.Args()
    if len(args) < 1 {
		usage()
	}
    
    util.OnInterrupt(func() {
		exitStatus=2
        glog.Flush()
		os.Exit(2)
	})
    if args[0] == "help" {
		help(args[1:])
		for _, cmd := range commands {
			if len(args) >= 2 && cmd.Name() == args[1] && cmd.ExecuteFunc != nil {
				fmt.Fprintf(os.Stderr, "Default Parameters:\n")
				cmd.Flag.PrintDefaults()
			}
		}
		return
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] && cmd.ExecuteFunc != nil {
			cmd.Flag.Usage = func() { cmd.Usage() }
			cmd.Flag.Parse(args[1:])
			args = cmd.Flag.Args()
			if !cmd.ExecuteFunc(cmd, args) {
				fmt.Fprintf(os.Stderr, "\n")
				cmd.Flag.Usage()
				fmt.Fprintf(os.Stderr, "Default Parameters:\n")
				cmd.Flag.PrintDefaults()
			}
			return
		}
	}
    util.Debug("start KMS at port:3000!")
	
}
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	//t.Funcs(template.FuncMap{"trim": strings.TrimSpace, "capitalize": capitalize})
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}
func usage() {
    printUsage(os.Stderr)
	fmt.Fprintf(os.Stderr, "For Logging, use \"kms [logging_options] [command]\".")
    flag.PrintDefaults()
	os.Exit(2)
}
func printUsage(w io.Writer) {
	tmpl(w, usageTemplate, commands)
}
var helpTemplate = `Usage: kms {{.UsageLine}}
  {{.Long}}
`
func help(args []string) {
	if len(args) == 0 {
		printUsage(os.Stdout)
		// not exit 2: succeeded at 'kms help'.
		return
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: kms help command\n\nToo many arguments given.\n")
		os.Exit(2) // failed at 'kms help'
	}

	arg := args[0]

	for _, cmd := range commands {
		if cmd.Name() == arg {
			tmpl(os.Stdout, helpTemplate, cmd)
			// not exit 2: succeeded at 'kms help cmd'.
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown help topic %#q.  Run 'kms help'.\n", arg)
	os.Exit(2) // failed at 'kms help cmd'
}