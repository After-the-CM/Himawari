package scanner

import (
	"Himawari/models/entity"
	"Himawari/models/logger"
	"net/http"
	"net/http/cookiejar"
)

var jar, _ = cookiejar.New(nil)
var client = &http.Client{
	Jar: jar,
	/*
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	*/
	Transport: logger.LoggingRoundTripper{
		Proxied: http.DefaultTransport,
	},
}

func Scan(j *entity.JsonNode) {

	Osci(j)
	if len(j.Children) > 0 {
		for i := 0; i < len(j.Children); i++ {
			Scan(&j.Children[i])
		}
	}

}
