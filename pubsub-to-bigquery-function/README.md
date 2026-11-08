# Pub/Sub to BigQuery Cloud Function (Go) & Blackbox Contract Testing

Google Cloud Pub/Sub のメッセージを受信し、**BigQuery** にレコードをストリーミング挿入する Cloud Functions と、**WireMock × Hurl × mise** による言語非依存ブラックボックステスト環境です。

BigQuery エミュレーターを起動するのではなく、BigQuery クライアントが裏側で発行する HTTP REST API (`tabledata.insertAll`) を **WireMock** で受け止め、投げられた JSON ペイロードを **Hurl** で検証します。

---

## 📁 ディレクトリ構成

```text
pubsub-to-bigquery-function/
├── .github/
│   └── workflows/
│       └── test.yml                 # GitHub Actions (Lint + Blackbox Test)
├── src/                             # ★ アプリケーション本体
│   ├── go.mod                       # Go 1.27 + cloud.google.com/go/bigquery
│   └── function.go                  # BigQuery Streaming Insert 最小構成
├── tests/                           # ★ テスト環境一式 (自作モックコードゼロ)
│   ├── compose.blackbox.yaml        # テスト専用 Compose (WireMock + Function)
│   ├── wiremock/
│   │   └── mappings/
│   │       └── bigquery_insert_all.json  # BigQuery API ダミー応答定義
│   └── cases/
│       └── 01_insert_row.hurl       # Hurl テストケース
├── .golangci.yml                    # golangci-lint 設定ファイル
├── mise.toml                        # ツール & タスク定義
└── README.md
```

---

## 🛠️ テスト実行

```bash
mise run test
```
