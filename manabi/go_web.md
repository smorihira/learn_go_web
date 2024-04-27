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
        user, _ := data.UserByEmail(r.PostFormValue("email")) // 中でDBアクセスするのか？
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
    - {{ .Uuid }}のようにデータを埋め込みことができる
      - データはどこから？`ExecuteTemplate`に渡したthreadsと関係ある？
      


- いろいろ
    - `go run [ファイル名] &`でバックグラウンド実行できる
