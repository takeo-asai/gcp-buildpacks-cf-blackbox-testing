# Pub/Sub to HTTP Cloud Function (Go) & Blackbox Contract Testing

Google Cloud Pub/Sub のメッセージを受信し、外部エンドポイントに HTTP GET パラメータとして転送する Cloud Functions と、**WireMock × Hurl × mise** による言語非依存ブラックボックステスト環境です。

---

## 📁 ディレクトリ構成

```text
pubsub-to-http-function/
├── .github/
│   └── workflows/
│       └── test.yml                 # GitHub Actions (Lint + Blackbox Test)
├── src/                             # ★ アプリケーション本体
│   ├── go.mod                       # Go 1.27
│   └── function.go                  # 最小構成の Cloud Function
├── tests/                           # ★ テスト環境一式 (自作コードゼロ)
│   ├── compose.blackbox.yaml        # テスト専用 Compose (WireMock + Function)
│   ├── wiremock/
│   │   └── mappings/
│   │       └── get_root.json        # WireMock ダミー応答定義
│   └── cases/                       # Hurl テストケース群 (1ケース = 1ファイル)
│       ├── 01_normal_string.hurl
│       ├── 02_multibyte_urlencode.hurl
│       └── 03_json_payload.hurl
├── .golangci.yml                    # golangci-lint 設定ファイル
├── mise.toml                        # ツール (pack, hurl, go, golangci-lint) & タスク定義
└── README.md
```

---

## 🛠️ 主なコマンド (mise タスク)

```bash
# 1. コードのフォーマット
mise run fmt

# 2. 静的解析・リントチェック (golangci-lint)
mise run lint

# 3. コンテナビルド & ブラックボックステスト実行
mise run test

# 4. 上記すべてを一括実行
mise run check
```
