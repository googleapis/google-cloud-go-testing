// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	// This example demonstrates building a simple mock that counts the number
	// of Bucket calls before calling the real client and returning its output.

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
