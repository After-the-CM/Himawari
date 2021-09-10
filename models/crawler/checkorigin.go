package crawler

import (
	"fmt"
	"net/url"
	"os"

	"Himawari/models/entity"
)

const (
	empty     = ""
	httpPort  = "80"
	httpsPort = "443"
	httpSch   = "http"
	httpsSch  = "https"
)

func IsSameOrigin(t *entity.RequestStruct, n *url.URL) bool {
	fmt.Println(&n.Scheme)
	fmt.Printf("%p\n", &t.Referer)

	switch {
	//2つのポート番号が空の場合→ホスト(ポート番号を除く)と、スキームを比較する
	case (t.Referer.Port() == empty) && (n.Port() == empty):
		if (t.Referer.Hostname() == n.Hostname()) && (t.Referer.Scheme == n.Scheme) {
			//fmt.Println("true : 2つのポート番号が空の場合→ホスト(ポート番号を除く)と、スキームを比較する")
			return true
		} else {
			//fmt.Println("false : 2つのポート番号が空の場合→ホスト(ポート番号を除く)と、スキームを比較する")
			return false
		}
		//ホスト(ポート番号を含む)とスキームを比較する。
	case t.Referer.Host == n.Host:
		if t.Referer.Scheme == n.Scheme {
			//fmt.Println("true : ホスト(ポート番号を含む)とスキームを比較する。")
			return true
		} else {
			//fmt.Println(t.Referer)
			//fmt.Println(t.Path)
			//fmt.Println("false : ホスト(ポート番号を含む)とスキームを比較する。")
			return false
		}
		//どちらか片方がポートが空の場合
	case t.Referer.Port() == empty:
		if getSchemaPort(&(t.Referer.Scheme), n.Port()) {
			//fmt.Println("true : t.Referer.Port() == empty:")
			return true
		}
		fallthrough
		//どちらか片方がポートが空の場合
	case n.Port() == empty:
		if getSchemaPort(&(n.Scheme), t.Referer.Port()) {
			//fmt.Println("true : t.Referer.Port() == empty:")
			return true
		}
		//fmt.Println("false : リンクを見つけたページと、リンクのスキームとポートが自動解決できません。")
		return false

	default:
		//fmt.Println(t.Referer)
		//fmt.Println(t.Path)
		//fmt.Println("false : default")
		return false

	}
}

//mapで`http`,`https`を受け取ったらポート番号を返す
func getSchemaPort(s *string, p string) bool {

	switch *s {
	case httpSch:
		return httpPort == p
	case httpsSch:
		return httpsPort == p
	default:
		fmt.Fprintln(os.Stderr, "http,httpsのスキーム以外のポートは自動解決されません。")
		//ありえないポート番号をリターンさせる
		return false
	}

}

/*
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

*/

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
