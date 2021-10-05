package logger

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

type LoggingRoundTripper struct {
	Proxied http.RoundTripper
}

func LoggingSetting() {
	layout := "2006-01-02_15:04:05"
	dirName := "log"
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		os.Mkdir(dirName, 0666)
	}
	t := time.Now()
	fileName := "log/" + t.Format(layout) + ".log"
	logFile, _ := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	log.SetFlags(log.Flags() &^ log.LstdFlags)
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.SetPrefix("======================================================\n")
}

func (lrt LoggingRoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	dumpedReq, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	log.SetFlags(log.Ltime)
	log.Println(string(req.URL.Scheme) + "://" + string(req.URL.Host))
	log.SetFlags(log.Flags() &^ log.LstdFlags)
	log.Println(string(dumpedReq))
	log.Printf("\n\n\n")

	res, e = lrt.Proxied.RoundTrip(req)

	if e != nil {
		fmt.Fprintln(os.Stderr, e)
	} else {
		dumpedResp, err := httputil.DumpResponse(res, true)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		log.SetFlags(log.Ltime)
		log.Println(string(res.Request.URL.Scheme) + "://" + string(res.Request.URL.Host))
		log.SetFlags(log.Flags() &^ log.LstdFlags)
		log.Println(string(dumpedResp))
		log.Printf("\n\n\n")
	}

	return
}
