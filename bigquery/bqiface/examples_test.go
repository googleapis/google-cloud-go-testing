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

package bqiface_test

import (
	"context"

	"cloud.google.com/go/bigquery"
	"github.com/googleapis/google-cloud-go-testing/bigquery/bqiface"
)

func Example_AdaptClient() {
	ctx := context.Background()
	c, err := bigquery.NewClient(ctx, "my-project")
	if err != nil {
		// TODO: Handle error.
	}
	client := bqiface.AdaptClient(c)
	defer client.Close()
	ds := client.Dataset("my_dataset")
	md, err := ds.Metadata(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	_ = md // TODO: use md.
}
