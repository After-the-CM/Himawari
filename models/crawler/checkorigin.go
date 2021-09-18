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

func IsSameOrigin(r *entity.RequestStruct, n *url.URL) bool {

	switch {
	//2つのポート番号が空の場合→ホスト(ポート番号を除く)と、スキームを比較する
	case (r.Referer.Port() == empty) && (n.Port() == empty):
		if (r.Referer.Hostname() == n.Hostname()) && (r.Referer.Scheme == n.Scheme) {
			return true
		} else {
			return false
		}
		//ホスト(ポート番号を含む)とスキームを比較する。
	case r.Referer.Host == n.Host:
		if r.Referer.Scheme == n.Scheme {
			return true
		} else {
			return false
		}
		//どちらか片方がポートが空の場合
	case r.Referer.Port() == empty:
		if getSchemaPort(&(r.Referer.Scheme), n.Port()) {
			return true
		}
		fallthrough
		//どちらか片方がポートが空の場合
	case n.Port() == empty:
		if getSchemaPort(&(n.Scheme), r.Referer.Port()) {
			return true
		}
		return false

	default:
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

func JudgeMethod(r *entity.RequestStruct) {
	if r.Form.Method == "GET" {
		GetRequest(r)
	} else if r.Form.Method == "POST" {
		PostRequest(r)
	} else {
		fmt.Fprintln(os.Stderr, "Other Methods.")
		return
	}
}
