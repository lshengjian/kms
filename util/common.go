package util
import (
    "fmt"
    "os"
   	"os/signal"
	"syscall"
	"encoding/json"
	"net/http"
    "golang.org/x/text/encoding/simplifiedchinese"
    "github.com/golang/glog"
)
const (
	VERSION = "0.1"
)
func WriteJson(w http.ResponseWriter, r *http.Request, httpStatus int, obj interface{}) (err error) {
	var bytes []byte
	if r.FormValue("pretty") != "" {
		bytes, err = json.MarshalIndent(obj, "", "  ")
	} else {
		bytes, err = json.Marshal(obj)
	}
	if err != nil {
		return
	}
	callback := r.FormValue("callback")
	if callback == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpStatus)
		_, err = w.Write(bytes)
	} else {
		w.Header().Set("Content-Type", "application/javascript")
		w.WriteHeader(httpStatus)
		if _, err = w.Write([]uint8(callback)); err != nil {
			return
		}
		if _, err = w.Write([]uint8("(")); err != nil {
			return
		}
		fmt.Fprint(w, string(bytes))
		if _, err = w.Write([]uint8(")")); err != nil {
			return
		}
	}

	return
}

func Utf8ToGBK(text string) (string, error) {
    dst := make([]byte, len(text)*2)
    tr := simplifiedchinese.GB18030.NewEncoder()
    nDst, _, err := tr.Transform(dst, []byte(text), true)
    if err != nil {
        return text, err
    }
    return string(dst[:nDst]), nil
}
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
		}
       	os.Exit(0)
	}()
}

func Debug(params ...interface{}) {
	glog.V(4).Infoln(params)
}
func Info(params ...interface{}) {
	glog.V(3).Infoln(params)
}
