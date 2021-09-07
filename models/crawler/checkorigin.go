package crawler

import (
	"net/url"

	"Himawari/models/entity"
)

const (
	empty     = ""
	httpPort  = "80"
	httpsPort = "443"
	httpSch   = "http"
	httpsSch  = "https"
)

/*
func main() {
	test := testStruct{"http://amazon/", "http://amazon:80/index.php"}
	fmt.Println(checkUrlOrigin(&test))
}
*/

func CheckUrlOrigin(t *entity.TestStruct) bool {
	base, _ := url.Parse(t.Origin)
	add, _ := url.Parse(t.Validation)
	//fmt.Println(base.Port())
	//fmt.Println(add.Port())

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

func JudgeMethod(r entity.RequestStruct) {
	if r.Form.Method == "GET" {
		GetRequest(r)
	} else if r.Form.Method == "POST" {
		PostRequest(r)
	} else {
		fmt.Println("Other Methods.")
		return 
	}
}
