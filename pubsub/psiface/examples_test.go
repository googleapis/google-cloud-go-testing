// Copyright 2019 Google LLC
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

package psiface_test

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/googleapis/google-cloud-go-testing/pubsub/psiface"
)

func ExampleAdaptClient() {
	ctx := context.Background()
	c, err := pubsub.NewClient(ctx, "")
	if err != nil {
		// TODO: Handle error.
	}
	client := psiface.AdaptClient(c)
	msg := psiface.AdaptMessage(&pubsub.Message{})
	_, err = client.Topic("my-topic").Publish(ctx, msg).Get(ctx)
	if err != nil {
		// TODO: Handle error.
	}
}
