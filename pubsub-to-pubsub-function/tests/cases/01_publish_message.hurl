# ケース 1: 受信した Pub/Sub メッセージが下流の Pub/Sub トピックに gRPC Publish されるか検証

# 1. Cloud Function に Pub/Sub イベント送信 ("hello-downstream-pubsub" -> base64: "aGVsbG8tZG93bnN0cmVhbS1wdWJzdWI=")
POST http://127.0.0.1:8080/
Content-Type: application/json
ce-id: event-001
ce-specversion: 1.0
ce-type: google.cloud.pubsub.topic.v1.messagePublished
ce-source: //pubsub.googleapis.com/projects/my-proj/topics/upstream-topic
[Options]
retry: 10
retry-interval: 500
{
    "message": {
        "data": "aGVsbG8tZG93bnN0cmVhbS1wdWJzdWI=",
        "messageId": "msg-ps-001"
    }
}
HTTP 200

# 2. GripMock の HTTP 管理 API (/api/history) を叩き、gRPC で届いた Publish リクエストを JSON として検証
GET http://127.0.0.1:4771/api/history
[Options]
retry: 5
retry-interval: 500
HTTP 200
[Asserts]
jsonpath "$[0].service" == "google.pubsub.v1.Publisher"
jsonpath "$[0].method" == "Publish"
jsonpath "$[0].requests[0].topic" == "projects/my-project/topics/downstream-topic"
# Protobuf の bytes フィールドは JSON 表現では base64 文字列として格納される
jsonpath "$[0].requests[0].messages[0].data" == "aGVsbG8tZG93bnN0cmVhbS1wdWJzdWI="
