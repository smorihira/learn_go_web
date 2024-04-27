- Go modules
    - Goのパッケージ管理ツール
    - 次のようにプロジェクトを作成する
        - プロジェクトのルートディレクトリを作成し `go mod init [ディレクトリ名]` を打つ
            - `go.mod` というGemFileやpackage.json的なものができる
        - ルートディレクトリ直下に　`main.go` を置けばプログラムを実行できる
        - `go run main.go` でビルド+実行できる　( `go build main.go` は実行はせずビルドのみ行う)
    - パッケージについて
        - `go get [パッケージ名]` で追加( `yarn add` と同じ)
        - `go mod tidy` で必要なor不要なパッケージを自動で追加or削除する
        - 同ディレクトリなら別ファイルでも同じパッケージとみなす(パッケージ構造=ディレクトリ構造になる)
        
- net/http
    - アプリケーションサーバを提供
        - `rails server` のようなサーバを立ち上げる操作は不要
        - そういう操作も全部含めてコンパイルされるため`go run`でビルドファイルを実行するだけでよい
    - マルチプレクサの提供
      ``` go
      mux := http.NewServeMux() // マルチプレクサを作成
      mux.HandleFunc([URL], [ハンドラ]) // ルーティング
      ```
        - マルチプレクサは静的ファイルを直接返送することもできる
          ``` go
          files := http.FileServer(http.Dir([path])) // 送信したいファイルを探すためのパスを指定
          mux.Handle([URL], http.StripPrefix([URL], files)) // ファイル取得のためのアクセスポイントの指定
          ```
        - これは[path]で始まるすべてのリクエストURLの[path]部分を[URL]に置き換えたpathにあるファイルを送信するコード

- クッキー・セッションの扱い
    - 認証ハンドラ(フォームを受け取り、認証してクッキーをセット)
      ``` go
      func authenticate(w http.ResponseWriter, r *http.Request) {
        r.ParseForm()
        user, _ := data.UserByEmail(r.PostFormValue("email")) // dataパッケージは後で作る
        if user.Password == data.Encrypt(r.PostFormValue("password")) { // ハッシュ化(偉い)
            session := user.CreateSession() // セッション作成(偉い)
            cookie := http.Cookie { // クッキーの作成(まあまあ)
                Name:       "_cookie",
                Value:      session.Uuid, // セッションIDをクッキーに保存
                HttpOnly:   true, // XSSとかでJSからクッキー取られないための設定？
            }
            http.SetCookie(w, &cookie) // ブラウザ上にクッキーを保存(偉い)
            http.Redirect(w, r, "/", 302)
        } else {
            http.Redirect(w, r, "/login", 302)
        }
      }
      ```
    - セッションは次のstructになっている
        - セッション情報はサーバ側に保存されており、ブラウザにはセッションidのみ保存される
      ``` go
      type Session {
            Id          int
            Uuid        string // セッションID
            Email       string
            UserId      int
            CreateAt    time.Time
      }
      ```
    - ログイン済みかチェック
        - ブラウザからクッキーを受け取り、その中のセッションIDと一致するセッションがサーバ内に存在するかチェックをする
        - ログイン済みかどうかでヘッダーを変えたりできる
      ``` go
      func session(w http.ResponseWriter, r *http.Request) (sess data.Session, err error) {
        cookie, err := r.Cookie("_cookie")
        if err != nil {
            return 
        }
        sess = data.Session {Uuid: cookie.Value}
        if ok, _ := sess.Check(); !ok {
            err = errors.New("Invalid session")
        }
        return
      }
      ```

- HTMLレスポンスの生成
    - 「アクション」と呼ばれる{{...}}で囲まれたコマンドが埋め込まれたhtmlファイルを用意する
    - その配列を`template.ParseFiles`でパースしてテンプレートを作成し、`ExecuteTemplate`メソッドでhtmlレスポンスを生成する
      ```go
      thread, err := data.Threads(); if err == nil { // これは何？
        // 本来はログイン済みかどうかで返すファイルが異なるため分岐する
        // _, err := session(w, r) 
        public_tmpl_files := []string{"templates/layout.html", "templates/public.navbar.html", "templates/index.html"}
        // パースでエラーを返す場合があるためMustを挟む必要がある
        var templates *template.Template
        templates = template.Must(template.ParseFiles(public_tmpl_files...))
        // トップのテンプレートファイルのみ実行すればよい
        templates.ExecuteTemplate(w, "layout", threads)
      }
      ```

- テンプレートファイルの中身について
    - 各テンプレートファイル(=アクションが埋め込まれたhtmlファイル)が一つのテンプレートを定義する
      ``` html
      {{ define "layout" }} <!-- テンプレート名 -->
      <!DOCTYPE html>
      <html lang="en">
        <head>
          <meta charset="utf-8">
          ...
        </head>
        <body>
          {{ template "navbar" . }}
          <div class="container">
            {{ template "content" . }}
          </div> <!-- container -->
          <script src="/static/js/jquery-2.1.1.min.js"></script>
          <script src="/static/js/bootstrap.min.js"></script>
        </body>
      </html>
      {{ end }}
      ```
    - defineアクションがテンプレートの宣言
    - templateアクションが他のテンプレートのレンダリング(railsでいうパーシャル)
        - テンプレート名とファイル名は一致する必要ない
        - 同じテンプレート名のものが複数ある場合はどうなるの？
    - 構造体のスライスを渡し{{ .Topic }}のようにデータを埋め込みことができる
        - データは`ExecuteTemplate`に渡す
        - {{ range . }}...{{ end }}で渡されたデータの数だけ繰り返しhtmlを生成できる
        - 与えられた構造体のフィールド(Topic等)だけでなく、その構造体のメソッドも{{ .NumReplies }}のように呼び出すことができる
          - rails的感覚だとビューの中でDBアクセスするのが気持ち悪い

- DBの扱い
    - postgre上でDBを作成
        - `createdb [DB名]`
    - 初期設定用に実行するクエリ(テーブル作成等)をまとめたスクリプトファイル(.sql)を用意し実行できる
        - `psql -f [sqlファイル名] -d [DB名]`
    - DBへの接続
        - DBへの接続プールを表すsql.DBへのポインタ型グローバル変数を用意し、この変数を介してクエリを発行する
        - プールを複数個用意してGoルーチンを用いて並行的にDBへ接続することも可能
     ``` go
     import "database/sql"
     var Db *sql.DB

     func init() {
     	var err error
     	Db, err = sql.Open("postgres", "dbname=chitchat sslmode=disable")
     	if err != nil {
     		log.Fatal(err)
     	}
     	return
     } 
     ```
    - クエリの発行は*sql.Rows型の値を返すQuery関数により次のように行える
        - `Db.Query("select ... where thread_id = $1", thread.Id)`
        - 返り値が1行のみな場合は*sql.Row型の値を返すQueryRow関数を用いることもできる

- サーバの起動
    - 構造体Serverを作成し、`ListenAndServe`を呼び出すだけで起動可能
      ``` go
      server := &http.Server {
        Addr:     "0.0.0.0:8000",
        Handler:  mux,
      }
      server.ListenAndServe()
      ```

- まとめ (`net/http`と`html/template`と`database/sql`がすごいだけ)
    - マルチプレクサを作成し`HandleFunc`で各ハンドラへルーティングをする
    - アクションが埋め込まれたテンプレートファイルを作成する
    - 各ハンドラ内で次の処理を行う
        - DBアクセスをしデータの取得等をする
        - テンプレートファイルを`ParseFiles`によりパースしテンプレートの作成
        - データとテンプレートを渡して`ExecuteTemplate`を実行する
    - ポートを指定してサーバを作成し`ListenAndServe`によりサーバを起動する
      


- いろいろ
    - `go run [ファイル名] &`でバックグラウンド実行できる
