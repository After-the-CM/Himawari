package crawler

import (
	"fmt"
	"net/http"
	"net/url"

	"Himawari/models/entity/entity.go"
)

func getRequest(r entity.RequestStruct) {

	u, _ := url.Parse(r.referer)
	a, _ := url.Parse(r.path)
	s := u.ResolveReference(a).String()

	Request, err := http.NewRequest("GET", s, nil)
	if err != nil {
		fmt.Println(err)
	}
	Request.URL.RawQuery = g.param.Encode()
	Request.Header.Set("User-Agent", "Himawari")

	client := new(http.Client)

	Response, err := client.Do(Request)

	defer Response.Body.Close()

	if err != nil {
		fmt.Println("Unable to reach the server.")
	} else {
		// bodyをfunc2に投げる。
		// func2(Response)
	}

	return
}
