package crawler

import (
	"fmt"
	"net/http"
	"net/url"

	"Himawari/models/entity"
)

func GetRequest(r entity.RequestStruct) {

	base, _ := url.Parse(r.Referer)
	rel, _ := url.Parse(r.Path)
	abs := base.ResolveReference(rel).String()

	Request, err := http.NewRequest("GET", abs, nil)
	if err != nil {
		fmt.Println(err)
	}

	client := new(http.Client)

	Response, err := client.Do(Request)

	Response.Body.Close()

	Request.URL.RawQuery = r.Param.Encode()
	Request.Header.Set("User-Agent", "Himawari")

	if err != nil {
		fmt.Println("Unable to reach the server.")
	} else {
		// bodyをfunc2に投げる。
		// func2(Response)
	}

	return
}
