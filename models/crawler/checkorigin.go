package crawler

import (
	"fmt"
	"net/url"
)

const (
	empty     = ""
	httpPort  = "80"
	httpsPort = "443"
	httpSch   = "http"
	httpsSch  = "https"
)

type testStruct struct {
	//リンクが存在したページのURL
	referer string
	//formの場合はaction
	path string
}

/*
func main() {
	test := testStruct{"http://amazon/", "http://amazon:80/index.php"}
	fmt.Println(checkUrlOrigin(&test))
}
*/

func checkUrlOrigin(t *testStruct) bool {
	base, _ := url.Parse(t.referer)
	add, _ := url.Parse(t.path)
	fmt.Println(base.Port())
	fmt.Println(add.Port())

	switch {
	//2つのポート番号が空の場合→ホスト(ポート番号を除く)と、スキームを比較する
	case (base.Port() == empty) && (add.Port() == empty):
		if (base.Hostname() == add.Hostname()) && (base.Scheme == add.Scheme) {
			return true
		} else {
			return false
		}
		//ホスト(ポート番号を含む)とスキームを比較する。
	case base.Host == add.Host:
		if base.Scheme == add.Scheme {
			return true
		} else {
			return false
		}
		//どちらか片方がポートが空の場合
	case base.Port() == empty:
		if (add.Port()) == (getSchemaPort(base.Scheme)) {
			return true
		}
		fallthrough
		//どちらか片方がポートが空の場合
	case add.Port() == empty:
		if (base.Port()) == (getSchemaPort(add.Scheme)) {
			return true
		}
		fallthrough

	default:
		return false

	}
}

//mapで`http`,`https`を受け取ったらポート番号を返す
func getSchemaPort(s string) string {
	ports := map[string]string{
		httpSch:  httpPort,
		httpsSch: httpsPort,
	}
	return ports[s]
}
