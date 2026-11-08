// Package function provides a minimal Cloud Function to insert Pub/Sub messages into BigQuery.
package function

import (
	"context"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"google.golang.org/api/option"
)

func init() {
	functions.CloudEvent("ProcessPubSubToBQ", ProcessPubSubToBQ)
}

// MessagePublishedData represents the data structure of a Pub/Sub CloudEvent.
type MessagePublishedData struct {
	Message struct {
		Data      []byte `json:"data"`
		MessageID string `json:"messageId"`
	} `json:"message"`
}

// Row represents a single record to be inserted into BigQuery.
type Row struct {
	MessageID string `bigquery:"message_id"`
	Payload   string `bigquery:"payload"`
}

// ProcessPubSubToBQ receives a Pub/Sub CloudEvent and inserts it into BigQuery.
func ProcessPubSubToBQ(ctx context.Context, e event.Event) error {
	var msg MessagePublishedData
	_ = e.DataAs(&msg)

	var opts []option.ClientOption
	// BIGQUERY_EMULATOR_HOST または BIGQUERY_ENDPOINT が指定されている場合はモックへ接続
	if endpoint := os.Getenv("BIGQUERY_ENDPOINT"); endpoint != "" {
		opts = append(opts, option.WithEndpoint(endpoint), option.WithoutAuthentication())
	}

	client, err := bigquery.NewClient(ctx, "my-project", opts...)
	if err != nil {
		return nil
	}
	defer func() { _ = client.Close() }()

	inserter := client.Dataset("my_dataset").Table("my_table").Inserter()
	_ = inserter.Put(ctx, []*Row{
		{
			MessageID: msg.Message.MessageID,
			Payload:   string(msg.Message.Data),
		},
	})

	return nil
}
