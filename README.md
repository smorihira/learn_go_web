# このレポジトリについて
ただのGo言語でのweb開発を勉強するためのplay ground

# 参考書籍
[Goプログラミング実践入門　標準ライブラリでゼロからWebアプリを作る](https://www.amazon.co.jp/Go%E3%83%97%E3%83%AD%E3%82%B0%E3%83%A9%E3%83%9F%E3%83%B3%E3%82%B0%E5%AE%9F%E8%B7%B5%E5%85%A5%E9%96%80-%E6%A8%99%E6%BA%96%E3%83%A9%E3%82%A4%E3%83%96%E3%83%A9%E3%83%AA%E3%81%A7%E3%82%BC%E3%83%AD%E3%81%8B%E3%82%89Web%E3%82%A2%E3%83%97%E3%83%AA%E3%82%92%E4%BD%9C%E3%82%8B-impress-top-gear/dp/4295000965)


# プロジェクト構成
    go_project             # ルートディレクトリ
    ├─ cmd                 # 各機能をまとめるディレクトリ
    │   └─ myapp           # サンプルアプリ1
    │       ├─ server      # 
    │       ├─ server.go   # 
    │       └─ tekitou.go  # 
    │       
    ├─ manabi              # 適当に学んだこと書く (落書き)
    |   ├─ go_iroiro.md    # Go自体についていろいろ
    |   ├─ web_dev.md      # Goでのweb開発一般 
    │   └─ go_web.md       # web用のもう少しテクニカルなこと
    │   
    ├─ pkg                 # パッケージ
    │   
    ├─ README.md           # ここ
    └─ go.mod              # パッケージ管理