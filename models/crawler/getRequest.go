package crawler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"Himawari/models/entity"
)

func GetRequest(r entity.RequestStruct) {

	base, _ := url.Parse(r.Referer)
	rel, _ := url.Parse(r.Path)
	abs := base.ResolveReference(rel).String()

	Request, err := http.NewRequest("GET", abs, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	Request.URL.RawQuery = r.Param.Encode()
	Request.Header.Set("User-Agent", "Himawari")

	client := new(http.Client)
	Response, err := client.Do(Request)

	body, _ := io.ReadAll(Response.Body)
	Response.Body.Close()

	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to reach the server.")
	} else {
		// bodyをfunc2に投げる。
		// func2(Response)
	}

	return
}
