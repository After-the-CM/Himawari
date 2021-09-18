package crawler

import (
	"io"
	"log"
	"net/url"
	"os"
	"time"

	"Himawari/models/entity"
)

func Crawl(url *url.URL) {
	loggingSetting()
	req := entity.RequestStruct{
		Referer: url,
		Path:    url,
		Param:   url.Query(),
	}
	GetRequest(&req)
}

func loggingSetting() {
	layout := "2006-01-02_15:04:05"
	dirName := "log"
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		os.Mkdir(dirName, 0666)
	}
	t := time.Now()
	fileName := "log/" + t.Format(layout) + ".log"
	logFile, _ := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	log.SetFlags(log.Ldate | log.Ltime)
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
}
