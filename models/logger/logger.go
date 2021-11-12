package logger

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"syscall"
	"time"
)

type LoggingRoundTripper struct {
	Proxied http.RoundTripper
}

var DumpedReq []byte

func LoggingSetting() {
	layout := "2006-01-02_15:04:05"
	dirName := "log"
	defaultUmask := syscall.Umask(0)
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err := os.Mkdir(dirName, 0777)
		if ErrHandle(err) {
			time.Sleep(time.Second * 5)
		}
	}
	t := time.Now()
	fileName := "log/" + t.Format(layout) + ".log"
	logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if ErrHandle(err) {
		time.Sleep(time.Second * 5)
	}
	syscall.Umask(defaultUmask)
	log.SetFlags(log.Flags() &^ log.LstdFlags)
	log.SetOutput(logFile)
	log.SetPrefix("======================================================\n")
}

func (lrt LoggingRoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	dumpedReq, err := httputil.DumpRequestOut(req, true)
	ErrHandle(err)

	log.SetFlags(log.Ltime)
	log.Println(req.URL.Scheme + "://" + req.URL.Host)
	log.SetFlags(log.Flags() &^ log.LstdFlags)
	log.Println(string(dumpedReq))
	DumpedReq = dumpedReq
	log.Printf("\n\n\n")

	res, e = lrt.Proxied.RoundTrip(req)

	if !ErrHandle(e) {
		dumpedResp, err := httputil.DumpResponse(res, true)
		ErrHandle(err)

		log.SetFlags(log.Ltime)
		log.Println(res.Request.URL.Scheme + "://" + res.Request.URL.Host)
		log.SetFlags(log.Flags() &^ log.LstdFlags)
		log.Println(string(dumpedResp))
		log.Printf("\n\n\n")
	}

	return
}

func ErrHandle(err error) bool {
	if err != nil {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			fmt.Fprintln(os.Stderr, file+":"+fmt.Sprint(line)+"\x1b[31;1m", err, "\x1b[0m")
		}
		return true
	}
	return false
}
