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
	Request.Header.Set("Referer", r.Referer)

	client := new(http.Client)
	Response, err := client.Do(Request)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to reach the server.")
	}

	body, _ := io.ReadAll(Response.Body)
	Response.Body.Close()
	// bodyをfunc2に投げる。
	// func2(bytes.NewBuffer(body), base)

	return
}
