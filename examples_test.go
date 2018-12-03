package googlecloudgotesting

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	"github.com/googleapis/google-cloud-go-testing/storage/stiface"
)

type RecordingClient struct {
	stiface.Client
	bucketCalls int
}

func (rc *RecordingClient) Bucket(name string) stiface.BucketHandle {
	rc.bucketCalls++
	return rc.Client.Bucket(name)
}

// We do not need to implement methods that we don't want to record - by default
// the embedded type will be used.

func Example_recordBuckets() {
	// This example demonstrates building a simple mock that the number of
	// Bucket calls before calling the real client and returning its output.

	ctx := context.Background()
	c, err := storage.NewClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	client := stiface.AdaptClient(c)
	recordingClient := RecordingClient{client, 0}

	recordingClient.Bucket("my-bucket-1")
	recordingClient.Bucket("my-bucket-2")
	recordingClient.Bucket("my-bucket-3")

	fmt.Println(recordingClient.bucketCalls)
	// Output: 3
}
