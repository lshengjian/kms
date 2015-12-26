package main

import (
    "fmt"
    "os"
   	"os/signal"
	"syscall"
    "flag"
    "sync"
    "net/http"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
    "github.com/golang/glog"
)
func OnInterrupt(fn func()) {
	// deal with control+c,etc
	signalChan := make(chan os.Signal, 1)
	// controlling terminal close, daemon not exit
	signal.Ignore(syscall.SIGHUP)
	signal.Notify(signalChan,
		os.Interrupt,
		os.Kill,
		syscall.SIGALRM,
		// syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		for _ = range signalChan {
			fn()
			os.Exit(0)
		}
	}()
}

func debug(params ...interface{}) {
	glog.V(4).Infoln(params)
}
var exitStatus = 0
var exitMu sync.Mutex

func setExitStatus(n int) {
	exitMu.Lock()
	if exitStatus < n {
		exitStatus = n
	}
	exitMu.Unlock()
}
func usage() {
	fmt.Fprintf(os.Stderr, "For Logging, use \"kms [logging_options] [command]\".")
	os.Exit(2)
}
var atexitFuncs []func()

func atexit(f func()) {
	atexitFuncs = append(atexitFuncs, f)
}

func exit() {
	for _, f := range atexitFuncs {
		f()
	}
	os.Exit(exitStatus)
}

func main() {
	//flag.Usage = usage
    flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "logs")
	//flag.Set("v", "4")
	flag.Parse()
    args := flag.Args()
    if len(args) < 1 {
		usage()
	}
    atexit(func() {
       glog.Flush()
    })
    
    OnInterrupt(func() {
		exitStatus=2
		exit()
	})
	e := echo.New()
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	e.Static("/", "public")
	e.Index("public/index.html")

	e.Favicon("public/favicon.ico")

	e.Post("/upload", upload)

	e.Get("/hello", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello!")
	})

	//print("start at port:3000!");
    debug("start at port:3000!")
    
	e.Run(":3000")
}
