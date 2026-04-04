# web-sqlx-sqlite-user

Go + Gin + sqlx + SQLite + Tailwind CSS v4 によるユーザー管理 Web アプリケーションです。

## 技術スタック

| 分類 | 技術 |
|------|------|
| 言語 | Go 1.22 以上 |
| Web フレームワーク | [Gin](https://github.com/gin-gonic/gin) |
| DB アクセス | [sqlx](https://github.com/jmoiron/sqlx) |
| データベース | SQLite ([modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite)) |
| CSS | [Tailwind CSS v4](https://tailwindcss.com/) |
| テンプレート | Go `html/template` |

## ディレクトリ構成

```
.
├── install/
│   ├── schema.sql      # テーブル定義
│   └── seed.sql        # 初期データ
├── public/assets/css/  # ビルド済み CSS（自動生成）
├── src/
│   ├── config/             # 環境変数・設定
│   ├── controllers/        # ハンドラー（ユーザー CRUD）
│   ├── db/                 # DB クエリ層
│   ├── lib/
│   │   ├── database/sqlite/ # DB 接続管理
│   │   └── template/        # テンプレートヘルパー
│   ├── style/
│   │   └── input.css       # Tailwind CSS エントリーポイント
│   ├── templates/          # HTML テンプレート
│   ├── main.go
│   └── router.go
├── go.mod
├── Makefile
└── package.json
```

## セットアップ & ビルド

### 前提条件

- Go 1.22 以上
- Node.js 20 以上 / npm

### 初回セットアップ

依存パッケージのインストール・CSS ビルド・アプリ起動を一括実行します。

```bash
make dev
```

内部では以下を順に実行します：

```bash
go mod tidy        # Go モジュールの依存解決
npm install        # Node.js 依存パッケージのインストール
npm run build      # Tailwind CSS のビルド
go run ./src       # アプリ起動
```

### Go バイナリのビルド

CSS をビルドしたうえで `web-sqlx-sqlite-user`（Windows では `.exe`）を生成します。

```bash
make binary
```

内部では以下を実行します：

```bash
npm run build                            # CSS ビルド
go build -o web-sqlx-sqlite-user ./src  # バイナリ生成
```

生成したバイナリを直接実行する場合：

```bash
./web-sqlx-sqlite-user        # Linux / macOS
.\web-sqlx-sqlite-user.exe    # Windows
```

### `go run` で起動（開発時）

```bash
make run
```

### CSS のみビルド

```bash
make build
```

### CSS ウォッチモード（開発時）

```bash
make watch
```

### 生成物の削除

```bash
make clean
```

バイナリ・SQLite ファイル・生成済み CSS を削除します。

---

## `make` を使わない場合

`make` が利用できない環境（Windows など）向けの手順です。

### 初回セットアップ & 起動

```bash
go mod tidy
npm install
npm run build
go run ./src
```

### Go バイナリのビルド & 実行

```bash
npm run build
go build -o web-sqlx-sqlite-user ./src

./web-sqlx-sqlite-user        # Linux / macOS
.\web-sqlx-sqlite-user.exe    # Windows
```

### `go run` で起動（開発時）

```bash
npm run build
go run ./src
```

### CSS のみビルド

```bash
npm run build
```

### CSS ウォッチモード（開発時）

```bash
npm run watch
```

### 生成物の削除

```bash
# Linux / macOS
rm -f web-sqlx-sqlite-user user.sqlite3 public/assets/css/style.css

# Windows (PowerShell)
Remove-Item -ErrorAction SilentlyContinue web-sqlx-sqlite-user.exe, user.sqlite3, public\assets\css\style.css
```

## 環境変数

`.env` ファイルをプロジェクトルートに作成して設定できます。

| 変数名 | デフォルト値 | 説明 |
|--------|-------------|------|
| `SERVER_PORT` | `8080` | サーバーのポート番号 |
| `DATABASE_PATH` | `./user.sqlite3` | SQLite ファイルパス |

## ライセンス

MIT
