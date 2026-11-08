# ケース 1: 通常の文字列が BigQuery のテーブルに行挿入されるか検証

# 1. WireMock の受信履歴をリセット
DELETE http://127.0.0.1:8081/__admin/requests
HTTP 200

# 2. Pub/Sub イベント送信 ("hello-bigquery-data")
POST http://127.0.0.1:8080/
Content-Type: application/json
ce-id: event-001
ce-specversion: 1.0
ce-type: google.cloud.pubsub.topic.v1.messagePublished
ce-source: //pubsub.googleapis.com/projects/my-proj/topics/my-topic
[Options]
retry: 10
retry-interval: 500
{
    "message": {
        "data": "aGVsbG8tYmlncXVlcnktZGF0YQ==",
        "messageId": "msg-bq-001"
    }
}
HTTP 200

# 3. WireMock に届いた BigQuery REST API リクエストを検証
GET http://127.0.0.1:8081/__admin/requests
[Options]
retry: 5
retry-interval: 500
HTTP 200
[Asserts]
jsonpath "$.requests" count == 1
jsonpath "$.requests[0].request.method" == "POST"
jsonpath "$.requests[0].request.url" == "/bigquery/v2/projects/my-project/datasets/my_dataset/tables/my_table/insertAll?alt=json&prettyPrint=false"
# BigQuery に送信された JSON ボディの中身を検証
jsonpath "$.requests[0].request.body" contains "\"message_id\":\"msg-bq-001\""
jsonpath "$.requests[0].request.body" contains "\"payload\":\"hello-bigquery-data\""
