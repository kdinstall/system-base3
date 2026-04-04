.PHONY: build watch run binary dev clean

## CSS ビルド（Tailwind v4）
build:
	npm run build

## CSS ウォッチ（バックグラウンド推奨）
watch:
	npm run watch

## Go アプリ起動（go run）
run: build
	go run ./src

## Go バイナリのビルド（CSS ビルドも含む）
binary: build
	go build -o web-sqlx-sqlite-user ./src

## 開発用: go mod tidy + npm install + CSS ビルド + 起動
dev:
	go mod tidy
	npm install
	npm run build
	go run ./src

## 生成物の削除
clean:
	rm -f web-sqlx-sqlite-user web-sqlx-sqlite-user.exe
	rm -f user.sqlite3 _user.sqlite3
	rm -f public/assets/css/style.css
