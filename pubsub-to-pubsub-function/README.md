# Pub/Sub to Pub/Sub Cloud Function (Go) & Blackbox Contract Testing

Google Cloud Pub/Sub のメッセージを受信し、下流の別 **Pub/Sub トピック** に Publish する Cloud Functions と、**GripMock × Hurl × mise** による言語非依存ブラックボックステスト環境です。

Pub/Sub エミュレーターを起動するのではなく、Pub/Sub クライアントが裏側で発行する gRPC (`google.pubsub.v1.Publisher/Publish`) を汎用 gRPC モックサーバー **GripMock** で受け止め、投げられた Protobuf ペイロードを HTTP 管理 API 経由で **Hurl** から JSON として検証します。

---

## 📁 ディレクトリ構成

```text
pubsub-to-pubsub-function/
├── .github/
│   └── workflows/
│       └── test.yml                 # GitHub Actions (Lint + Blackbox Test)
├── src/                             # ★ アプリケーション本体
│   ├── go.mod                       # Go 1.27 + cloud.google.com/go/pubsub
│   └── function.go                  # Pub/Sub Publish 最小構成
├── tests/                           # ★ テスト環境一式 (自作モックコードゼロ)
│   ├── compose.blackbox.yaml        # テスト専用 Compose (GripMock + Function)
│   ├── gripmock/
│   │   ├── proto/
│   │   │   └── pubsub.proto         # Pub/Sub Publisher 最小 proto 定義
│   │   └── stubs/
│   │       └── publish.yaml         # GripMock ダミー応答定義
│   └── cases/
│       └── 01_publish_message.hurl  # Hurl テストケース
├── .golangci.yml                    # golangci-lint 設定ファイル
├── mise.toml                        # ツール & タスク定義
└── README.md
```

---

## 🛠️ テスト実行

```bash
mise run test
```
