// Package function provides a minimal Cloud Function to forward a Pub/Sub message to another Pub/Sub topic.
package function

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

func init() {
	functions.CloudEvent("ProcessPubSubToPubSub", ProcessPubSubToPubSub)
}

// MessagePublishedData represents the data structure of a Pub/Sub CloudEvent.
type MessagePublishedData struct {
	Message struct {
		Data []byte `json:"data"`
	} `json:"message"`
}

// ProcessPubSubToPubSub receives a Pub/Sub CloudEvent and publishes its data to a downstream topic.
// 環境変数 PUBSUB_EMULATOR_HOST が設定されている場合、全言語の公式 SDK が自動的にモックサーバーへ平文 gRPC 接続します。
func ProcessPubSubToPubSub(ctx context.Context, e event.Event) error {
	var msg MessagePublishedData
	_ = e.DataAs(&msg)

	client, err := pubsub.NewClient(ctx, "my-project")
	if err != nil {
		return nil
	}
	defer func() { _ = client.Close() }()

	topic := client.Topic("downstream-topic")
	res := topic.Publish(ctx, &pubsub.Message{
		Data: msg.Message.Data,
	})
	_, _ = res.Get(ctx)

	return nil
}
