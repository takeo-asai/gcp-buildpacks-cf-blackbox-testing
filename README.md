# gcp-buildpacks-cf-blackbox-testing

Google Cloud Functions (2nd gen) を **Google Cloud Buildpacks (`pack`)** でコンテナ化し、言語内部の実装を一切知らなくても外部から振る舞いを検証できる **言語非依存のブラックボックステスト（コンテナ対向テスト）** のリファレンス実装です。

---

## 🎯 コンセプトと特徴

1. **完全な言語非依存 (Language-Agnostic)**
   - Cloud Functions (第2世代) の実体は「ポート `8080` で CloudEvents (HTTP POST) を待ち受ける OCI コンテナ」です。
   - `pack` CLI でビルドしたコンテナに対して外から HTTP リクエストを投げて検証するため、中身が **Go から Python / Node.js / Java / Ruby に置き換わっても、テストスイート（`tests/`）は 1 文字も変更せずにそのまま機能** します。

2. **エミュレーター不要・自作モックコードゼロ**
   - 重厚な GCP エミュレーターを起動したり、モックサーバーを自作するのではなく、通信境界（HTTP REST や gRPC）を実績ある汎用 OSS モックコンテナ（**WireMock** / **GripMock**）で受け止め、送信されたリクエストが仕様通りかを検証します。

3. **宣言的なテストとモダンなツール管理**
   - **Hurl**: リクエスト送信・ポーリング待機（`retry`）・JSON アサーションをプレーンテキストで宣言的に記述。
   - **Docker Compose (`compose.blackbox.yaml`)**: テスト対象コンテナとモックコンテナの配線を宣言的に定義。
   - **mise**: `pack`, `hurl`, `golangci-lint`, `go` のバージョン管理とタスク実行を `mise.toml` 1枚に統合。

---

## 🏗️ アーキテクチャ概要

```text
┌─────────────────────────────────────────────────────────────────┐
│ テスト実行環境 (Local / GitHub Actions CI)                       │
│                                                                 │
│   [Hurl (テストランナー)]                                        │
│     │                                                           │
│     ├─ ① CloudEvent (HTTP POST) を送信 ──▶ [Cloud Function]     │
│     │                                        (:8080 コンテナ)    │
│     │                                               │           │
│     │                                  ② 外部通信 (HTTP / gRPC) │
│     │                                               ▼           │
│     └─ ③ 受信履歴 API を叩いて検証 ──────▶ [汎用 OSS モック]    │
│          (JSONPath アサーション)           (WireMock / GripMock)│
└─────────────────────────────────────────────────────────────────┘
```

---

## 📂 収録サンプル一覧

本リポジトリには、Cloud Functions の代表的な 3 つの外部連携パターンのサンプルが含まれています。

| ディレクトリ | 処理内容 | 外部通信プロトコル | 使用モックコンテナ | 検証 API |
| :--- | :--- | :--- | :--- | :--- |
| **[1. pubsub-to-http-function](./pubsub-to-http-function/)** | Pub/Sub を受け取り、外部 Web API (`example.com`) へ GET 送信 | **HTTP GET** | **WireMock** (`wiremock/wiremock`) | `GET :8081/__admin/requests` |
| **[2. pubsub-to-bigquery-function](./pubsub-to-bigquery-function/)** | Pub/Sub を受け取り、**BigQuery** テーブルへ Streaming Insert | **HTTP POST** (REST JSON) | **WireMock** (`wiremock/wiremock`) | `GET :8081/__admin/requests` |
| **[3. pubsub-to-pubsub-function](./pubsub-to-pubsub-function/)** | Pub/Sub を受け取り、下流の **Pub/Sub** トピックへ Publish | **gRPC** (HTTP/2 + Protobuf) | **GripMock** (`bavix/gripmock`) | `GET :4771/api/history` |

---

## 📁 リポジトリ構成

```text
gcp-buildpacks-cf-blackbox-testing/
├── .github/
│   └── workflows/
│       └── ci.yml                       # GitHub Actions (Matrix で 3 プロジェクトを並列検証)
├── pubsub-to-http-function/             # パターン 1: 外部 HTTP GET 転送
│   ├── src/                             #   関数本体 (Go 1.27)
│   ├── tests/
│   │   ├── compose.blackbox.yaml        #   WireMock + Function 定義
│   │   ├── wiremock/mappings/           #   WireMock スタブ定義
│   │   └── cases/*.hurl                 #   Hurl テストケース (3ケース)
│   ├── .golangci.yml
│   └── mise.toml
├── pubsub-to-bigquery-function/         # パターン 2: BigQuery レコード挿入
│   ├── src/                             #   関数本体 (cloud.google.com/go/bigquery)
│   ├── tests/
│   │   ├── compose.blackbox.yaml        #   WireMock + Function 定義
│   │   ├── wiremock/mappings/           #   BigQuery insertAll スタブ定義
│   │   └── cases/*.hurl                 #   Hurl テストケース
│   ├── .golangci.yml
│   └── mise.toml
├── pubsub-to-pubsub-function/           # パターン 3: Pub/Sub トピック Publish
│   ├── src/                             #   関数本体 (cloud.google.com/go/pubsub)
│   ├── tests/
│   │   ├── compose.blackbox.yaml        #   GripMock + Function 定義
│   │   ├── gripmock/                    #   Pub/Sub proto & スタブ定義
│   │   └── cases/*.hurl                 #   Hurl テストケース
│   ├── .golangci.yml
│   └── mise.toml
├── mise.toml                            # ルートタスク定義
└── README.md
```

---

## 🚀 クイックスタート

### 前提条件
- **Docker** (および `docker compose`)
- **[mise](https://mise.jdx.dev/)** (`pack`, `hurl`, `go`, `golangci-lint` は `mise` が自動インストールします)

### 1. ツールのセットアップ
```bash
mise install
```

### 2. 全プロジェクトのテストを一括実行
```bash
mise run test
```

### 3. 個別のプロジェクトごとに実行する場合
```bash
# パターン 1: HTTP
mise run test:http

# パターン 2: BigQuery
mise run test:bigquery

# パターン 3: Pub/Sub (gRPC)
mise run test:pubsub
```

または各ディレクトリに移動して個別にタスクを実行できます：
```bash
cd pubsub-to-http-function
mise run fmt    # コード整形 (gofmt)
mise run lint   # 静的解析 (golangci-lint)
mise run build  # Buildpacks (pack) でコンテナビルド
mise run test   # Compose 起動 -> Hurl テスト -> Compose 停止
mise run check  # fmt + lint + test 一括実行
```
