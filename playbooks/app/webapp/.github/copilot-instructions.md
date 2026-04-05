# Copilot Instructions for webapp

## プロジェクト概要
- このリポジトリは Go + Gin + sqlx + SQLite + Tailwind CSS v4 で実装されたユーザー管理アプリです。
- エントリーポイントは `src/main.go`、ルーティングは `src/router.go` です。
- 画面は Go の `html/template` で `src/templates/*.html` を読み込みます。

## 技術と主要ディレクトリ
- Web: Gin (`github.com/gin-gonic/gin`)
- DB: sqlx + modernc sqlite (`github.com/jmoiron/sqlx`, `modernc.org/sqlite`)
- CSS: Tailwind CSS v4 (`@tailwindcss/cli`)
- 設定: `src/config/config.go`
- コントローラ: `src/controllers/`
- DB クエリ層: `src/db/`
- DB 接続基盤: `src/lib/database/sqlite/`
- テンプレート共通データ: `src/lib/template/`
- CSS 入力: `src/style/input.css`
- CSS 出力(生成物): `public/assets/css/style.css`

## 実装ルール
- 既存レイヤ構造を維持してください。
  - HTTP ハンドリングは `src/controllers/`
  - SQL/永続化ロジックは `src/db/`
  - DB 接続共通処理は `src/lib/database/sqlite/`
- 既存の命名規約とコメント言語（日本語コメント）を優先してください。
- DB 操作の追加・更新は、原則として `UserDb` と同様にトランザクションヘルパー `DoInTx` の利用を検討してください。
- テンプレートへ渡すデータは `tmpl.MergeData(...)` を使い、共通キー（`app_name`, `g_year`）を維持してください。
- 404 の描画は既存実装に合わせ、`404.html` を使用してください。

## 変更時の注意点
- `public/assets/css/style.css` はビルド生成物です。直接編集せず `src/style/input.css` を編集してください。
- SQLite の接続は `sync.OnceValue` で単一化されています。接続管理の方式を不用意に変更しないでください。
- 現在の環境変数名はコード実装を正とし、`SERVER_PORT` と `DATABASE_PATH` を使用してください。
- 既存仕様を変える場合は、関連するテンプレート・ルート・DB クエリを一貫して更新してください。

## 変更後の確認コマンド
- Go のビルド確認: `go build ./src`
- CSS ビルド: `npm run build`
- 開発起動: `go run ./src`
- Makefile を使う場合: `make build`, `make run`, `make binary`

## Copilot への期待動作
- 変更は最小差分で行い、無関係なリファクタリングは避けてください。
- 新規機能追加時は、ルーティング (`src/router.go`)・コントローラ・DB 層・テンプレートの整合を必ず取ってください。
- バリデーションエラー時の表示は既存パターン（`errors` 配列をテンプレートへ渡す）に合わせてください。
- SQL は既存スタイルに合わせてプレースホルダ `?` を使用してください。
- 可能であれば変更対象に応じて `go build ./src` または `npm run build` を実行し、結果を報告してください。
