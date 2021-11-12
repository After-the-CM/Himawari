package entity

type Issue struct {
	URL       string
	Kind      string
	Parameter string
	Payload   string
	Cookie    JsonCookie
	Request   string
	Response  string
}

type Vuln struct {
	CWE         string  `json:"CWE"`
	Severity    string  `json:"Severity"`
	Description string  `json:"Description"`
	Mandatory   string  `json:"Mandatory"`
	Insurance   string  `json:"Insurance"`
	Issues      []Issue `json:"Issues"`
}

var csrfVuln = Vuln{
	CWE:      "CWE-352",
	Severity: "High",
	Description: "クロスサイト・リクエスト・フォージェリ(CSRF)は、ログインした利用者の意図したリクエストであるかの確認ができていないと発生する脆弱性で、CSRFの脆弱性が存在するとアカウントの不正利用、アカウント内情報の改ざんなどが発生する可能性があります。\n" +
		"CSRFの特徴は、悪意のあるリクエストがユーザのブラウザを経由してWebアプリケーションに送られるため、正規のユーザが実際に許可したリクエストと見分けがつかないことです。\n",
	Mandatory: "対策の必要な部分にワンタイムトークンの生成とチェック機能を実装しましょう。トークンによってリクエストの正当性を確認することができます。\n" +
		"登録・確定等の重要な処理の前にはパスワードの再確認を行うことで正規のユーザであることを確認できます。\n" +
		"正規のリクエストとCSRF攻撃のリクエストではリクエストされているページへのリンク先を持った直前のURLが異なるので、Refererヘッダの確認は行いましょう。\n",
	Insurance: "新規デバイスのログイン、パスワード変更などの重要な操作を行った際は、アカウントに登録されているメールアドレスに確認を行いましょう。\n",
	Issues:    []Issue{},
}

var reflectedxssVuln = Vuln{
	CWE:      "CWE-79",
	Severity: "High",
	Description: "反射型クロスサイト・スクリプティング(XSS)は、攻撃者によって仕掛けられた悪意のあるスクリプトをサイトに訪れたユーザのブラウザ上で実行させることができる脆弱性です。\n" +
		"XSSの攻撃を受けたユーザは、サイトの利用者の権限でWebアプリケーションの機能を悪用されたり、Cookie値を盗まれて成りすましの被害にあう可能性があります。\n",
	Mandatory: "要素内容・属性値の文字をエスケープして、属性値はダブルクォーテーションでくくりましょう。\n" +
		"(例. `<` を `&lt;` にする等、HTML特殊文字はHTMLエンティティに置き換えましょう。)\n" +
		"HTTPレスポンスヘッダのContent-Typeフィールドに文字コードを指定しましょう。文字コードの指定を省略している場合、Webアプリケーションとブラウザとの文字エンコーディングの差異でXSSの原因になりえます。\n",
	Insurance: "X-XSS-Protectionレスポンスヘッダを利用することでブラウザ機能により反射型XSSを無害な出力に変更できます。X-XSS-Protectionレスポンスヘッダを利用していれば利用者が無効化している場合でも上書きして有効にすることができます。\n" +
		"CSP(Contents Security Policy)の設定することでevalを禁止、インラインスクリプトを禁止にする等、スクリプトに対してより強い制約を設けることができます。\n" +
		"入力値の検証は行ったほうがよいです。この対策は限定的ですが、郵便番号など文字の種類や入力値の長さを制限できる箇所に要求を満たしていない入力がされた場合はエラーを表示し、再入力を求めるようにすることでスクリプトを注入できないようにすることができます。\n" +
		"上記の検証はサーバ側で行ったほうがよいです。クライアント上で入力値の検証を行ったとしても、プロキシツールなどの使用により攻撃者は要求を満たしていない値の送信が可能です。これを防ぐために、サーバ側で適切な検証を行いましょう。\n" +
		"CookieにHttpOnly属性を付与するとCookieに保存されているセッションIDへのアクセスを防ぐことができます。",
	Issues: []Issue{},
}

var storedxssVuln = Vuln{
	CWE:      "CWE-79",
	Severity: "High",
	Description: "持続型クロスサイト・スクリプティング(XSS)は、攻撃者によって仕掛けられた悪意のあるスクリプトをサイトに訪れたユーザのブラウザ上で実行させることができる脆弱性です。\n" +
		"持続型XSSの脆弱性が存在すると、攻撃者によって書き込まれた悪意のあるスクリプトがサイトのDB等に保存されます。攻撃用のスクリプトがDBに保存されていると持続型XSSの脆弱性が存在したページにアクセスするたびにスクリプトが実行されてしまいます。\n",
	Mandatory: "要素内容/属性値の文字をエスケープして、属性値はダブルクォーテーションでくくりましょう。\n" +
		"(例. `<` を `&lt;` にする等、HTML特殊文字はHTMLエンティティに置き換えましょう。)\n" +
		"HTTPレスポンスヘッダのContent-Typeフィールドに文字コードを指定しましょう。文字コードの指定を省略している場合、Webアプリケーションとブラウザとの文字エンコーディングの差異でXSSの原因になりえます。\n",
	Insurance: "CSP(Contents Security Policy)の設定することでevalを禁止、インラインスクリプトを禁止にする等、スクリプトに対してより強い制約を設けることができます。\n" +
		"入力値の検証は行ったほうがよいです。この対策は限定的ですが、郵便番号など文字の種類や入力値の長さを制限できる箇所に要求を満たしていない入力がされた場合はエラーを表示し、再入力を求めるようにすることでスクリプトを注入できないようにすることができます。\n" +
		"上記の検証はサーバ側で行ったほうがよいです。クライアント上で入力値の検証を行ったとしても、プロキシツールなどの使用により攻撃者は要求を満たしていない値の送信が可能です。これを防ぐために、サーバ側で適切な検証を行いましょう。\n" +
		"CookieにHttpOnly属性を付与するとCookieに保存されているセッションIDへのアクセスを防ぐことができます。",
	Issues: []Issue{},
}

var osciVuln = Vuln{
	CWE:      "CWE-78",
	Severity: "High",
	Description: "OSコマンド・インジェクションは、シェルの不適切な呼び出し方をしている場合に意図しないOSコマンドの実行が可能になる脆弱性です。\n" +
		"OSコマンド・インジェクションの脆弱性が存在すると、攻撃者によって情報漏洩、任意のOSコマンドの実行、不正なシステム操作、他システムへの攻撃の踏み台などの攻撃を受ける可能性があります。\n",
	Mandatory: "可能な限り、シェルを呼び出す機能のある関数の利用は避けましょう。\n" +
		"ライブラリを使った実装に切り替えることができないかを検討してください。\n" +
		"シェルを呼び出す機能のある関数を利用する場合は、外部からのパラメータを渡さないように実装しましょう。\n",
	Insurance: "シェルを呼び出す機能のある関数を利用する場合は、その引数を構成する変数を調べ、許可した処理のみ実行するように実装しましょう。\n",
	Issues:    []Issue{},
}

var dirtraversalVuln = Vuln{
	CWE:         "CWE-22",
	Severity:    "High",
	Description: "ディレクトリ・トラバーサルは、攻撃者によって制限されたディレクトリ外の任意のファイルに対して閲覧をはじめとする、開発者の意図しない処理が行える脆弱性です。\n",
	Mandatory: "Webサーバ内のファイル名を外部からのパラメータで指定する実装は避けてください。\n" +
		"ファイルを開く部分の実装は、固定のディレクトリを指定して、かつファイル名にディレクトリ名が含まれないようにしましょう。\n",
	Insurance: "Webサーバ内のファイルへのアクセス権限の設定を正しく管理できていれば、Webサーバ側でアクセスを拒否することができます。\n" +
		"ファイル名の文字の種類を英数字に限定しておけば、ディレクトリ・トラバーサル攻撃時に利用される記号文字が使用できなくなります。\n",

	Issues: []Issue{},
}

var timebasedsqliVuln = Vuln{
	CWE:      "CWE-89",
	Severity: "High",
	Description: "ブラインドSQL インジェクションは、DBと連携したアプリケーションでSQLの呼び出し方に不備があると発生する脆弱性です。\n" +
		"通常のSQL インジェクションの脆弱性は脆弱性は存在するが情報漏洩させるのは難しい場面がありますが、ブラインドSQLインジェクションの攻撃手法で情報漏洩を引き起こすことが可能になります。情報漏洩の他にデータベースの改ざん、不正ログイン、OSコマンドの実行、ファイルの参照・更新などの攻撃が発生する可能性があります。\n",
	Mandatory: "SQL文の組み立ては静的プレースホルダで行うよう実装しましょう。\n" +
		"Webアプリケーションに渡されるパラメータでSQL文を直接指定しないようにしましょう。\n" +
		"文字列リテラル内で特別な意味を持つ記号文字「'」を「\\'」、「\\」を「\\\\」などに置き換えるエスケープ処理を行いましょう。\n",
	Insurance: "Webアプリケーションが利用するDBアカウントには適切な権限を与えてください。書き込み権限を与えなければ、情報の改ざんを防ぐことができます。\n",
	Issues:    []Issue{},
}

var errbasedsqliVuln = Vuln{
	CWE:      "CWE-89",
	Severity: "High",
	Description: "SQL インジェクションは、DBと連携したアプリケーションでSQLの呼び出し方に不備があると発生する脆弱性です。\n" +
		"SQL インジェクションの脆弱性が存在すると、攻撃者によって情報漏洩、データベースの改ざん、不正ログイン、OSコマンドの実行、ファイルの参照・更新などの攻撃が発生する可能性があります。\n",
	Mandatory: "SQL文の組み立ては静的プレースホルダで行うよう実装しましょう。\n" +
		"Webアプリケーションに渡されるパラメータでSQL文を直接指定しないようにしましょう。\n" +
		"文字列リテラル内で特別な意味を持つ記号文字「'」を「\\'」、「\\」を「\\\\」などに置き換えるエスケープ処理を行いましょう。\n",
	Insurance: "エラーメッセージがそのままブラウザに表示されないよう設定を行ってください。\n" +
		"Webアプリケーションが利用するDBアカウントには適切な権限を与えてください。書き込み権限を与えなければ、情報の改ざんを防ぐことができます。\n",
	Issues: []Issue{},
}

var openredirectVuln = Vuln{
	CWE:      "CWE-601",
	Severity: "Medium",
	Description: "オープンリダイレクトは、攻撃者が任意の外部ドメインにリダイレクトさせるURLを作成できる脆弱性です。\n" +
		"オープンリダイレクトの脆弱性が存在する場合、攻撃者は被害者を悪意のあるサイトに誘導することで、フィッシング攻撃が可能になります。\n" +
		"この攻撃は、正しいURLを使用している上にSSL証明書のエラーも出ないため被害者が気づきにくくなってしまっており、ログインIDとパスワードなどの個人情報を求められた場合に被害者は入力してしまう可能性があります。\n",
	Mandatory: "リダイレクト先のURLは固定するようにしましょう。\n" +
		"許可されたドメインのみにリダイレクトができる制限をかけましょう。\n",
	Issues: []Issue{},
}

var httpheaderiVuln = Vuln{
	CWE:      "CWE-113",
	Severity: "Medium",
	Description: "HTTP ヘッダインジェクションは、HTTPレスポンスヘッダの出力処理に問題がある場合に発生する脆弱性です。\n" +
		"HTTP ヘッダインジェクションの脆弱性が存在すると、成りすまし、表示内容の改変、キャッシュ汚染、任意のJavaScriptの実行などの攻撃が発生する可能性があります。\n",
	Mandatory: "ヘッダ出力用のライブラリやAPIを利用するようにしましょう。\n" +
		"ヘッダ生成を行うパラメータに対して改行を許可しないような処理を行うようにしましょう。\n" +
		"外部からのパラメータをHTTPレスポンスヘッダとして出力しないようにしましょう。\n",
	Insurance: "外部からの入力の改行コードはすべて削除するようにしましょう。改行コードを含む文字列を入力された場合はWebアプリケーションが正しく動作しなくなる場合があります。\n",
	Issues:    []Issue{},
}

var dirlistingVuln = Vuln{
	CWE:      "CWE-548",
	Severity: "Low",
	Description: "ディレクトリ・リスティングとは、Webサーバのディレクトリ内容をリスト表示するものです。特定のWebサイトディレクトリにインデックスファイルがない場合にディレクトリの内容を表示するWebサーバーの機能がオンになっていることで、ディレクトリ・リスティングが発生します。\n" +
		"この機能がオンになっていることによって意図せずファイルが公開され、攻撃者が攻撃を仕掛けるのに十分な情報を与えてしまう可能性があります。\n",
	Mandatory: "公開ディレクトリには非公開のファイルを配置しないようにしてください。ディレクトリ・リスティング機能によって外部から閲覧することが可能になり、情報漏洩となります。\n",
	Insurance: "Webサーバのディレクトリ・リスティング機能をオフにしてください。\n",
	Issues:    []Issue{},
}

var Vulnmap = map[string]*Vuln{
	"Cross_Site_Request_Forgery": &csrfVuln,
	"Reflected_XSS":              &reflectedxssVuln,
	"Stored_XSS":                 &storedxssVuln,
	"OS_Command_Injection":       &osciVuln,
	"Directory_Traversal":        &dirtraversalVuln,
	"Time_based_SQL_Injection":   &timebasedsqliVuln,
	"Error_Based_SQL_Injection":  &errbasedsqliVuln,
	"Open_Redirect":              &openredirectVuln,
	"HTTP_Header_Injection":      &httpheaderiVuln,
	"Directory_Listing":          &dirlistingVuln,
}

var WholeIssue []Issue
