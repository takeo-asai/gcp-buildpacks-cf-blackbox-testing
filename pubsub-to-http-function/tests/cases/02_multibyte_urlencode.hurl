# ケース 2: 日本語・記号（URLエンコードが必要な文字列）

# 1. WireMock の受信履歴をリセット
DELETE http://127.0.0.1:8081/__admin/requests
HTTP 200

# 2. Pub/Sub イベント送信 ("日本語&記号=100%")
POST http://127.0.0.1:8080/
Content-Type: application/json
ce-id: event-002
ce-specversion: 1.0
ce-type: google.cloud.pubsub.topic.v1.messagePublished
ce-source: //pubsub.googleapis.com/projects/my-proj/topics/my-topic
{
    "message": {
        "data": "5pel5pys6KqeJuiomOWPtj0xMDAl",
        "messageId": "msg-002"
    }
}
HTTP 200

# 3. WireMock 管理 API で受信リクエストを検証
GET http://127.0.0.1:8081/__admin/requests
[Options]
retry: 5
retry-interval: 500
HTTP 200
[Asserts]
jsonpath "$.requests" count == 1
jsonpath "$.requests[0].request.method" == "GET"
jsonpath "$.requests[0].request.queryParams.param.values[0]" == "日本語&記号=100%"
