# Example for Toramaru

Toramaruの簡単な使い方―API開発中のCORS回避

## 構成

- API Server: api.go
- Static File: index.html
- Proxy Server: toramaru

## 実験手順

1. サーバを起動

```bash
go run api.go &
python -m http.serve 8070 &
toramaru -p 8080 -r "/api/>localhost:8071" -r "/>localhost:8070"
```

2. ブラウザで http://localhost:8080 を開く
3. 文字列を入力してボタンを押す
4. 開発ツールでコンソールを開く
5. CORSが原因のエラーが出ていないことを確認
