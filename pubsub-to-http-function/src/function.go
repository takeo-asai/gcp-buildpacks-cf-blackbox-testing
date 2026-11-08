// Package function provides a minimal Cloud Function to forward Pub/Sub messages via HTTP GET.
package function

import (
	"context"
	"net/http"
	"net/url"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

func init() {
	functions.CloudEvent("ProcessPubSub", ProcessPubSub)
}

// MessagePublishedData represents the data structure of a Pub/Sub CloudEvent.
type MessagePublishedData struct {
	Message struct {
		Data []byte `json:"data"`
	} `json:"message"`
}

// ProcessPubSub receives a Pub/Sub CloudEvent and sends its data as a GET query parameter.
func ProcessPubSub(ctx context.Context, e event.Event) error {
	var msg MessagePublishedData
	_ = e.DataAs(&msg)

	baseURL := os.Getenv("TARGET_ENDPOINT")
	if baseURL == "" {
		baseURL = "https://example.com/"
	}

	endpoint := baseURL + "?param=" + url.QueryEscape(string(msg.Message.Data))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil
	}

	resp, _ := http.DefaultClient.Do(req)
	if resp != nil {
		_ = resp.Body.Close()
	}

	return nil
}
