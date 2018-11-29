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

package dsiface

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/datastore"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("integration tests skipped in short mode")
	}

	projID := os.Getenv("DATASTORE_PROJECT_ID")
	if projID == "" {
		t.Skip("missing DATASTORE_PROJECT_ID environment variable")
	}

	kind := fmt.Sprintf("dsiface_test_%d", time.Now().UnixNano())
	ctx := context.Background()
	c, err := datastore.NewClient(ctx, projID)
	if err != nil {
		t.Fatal(err)
	}
	client := AdaptClient(c)
	defer client.Close()
	basicTests(t, kind, client)
}

type sourceData struct {
	Payload string
}

func basicTests(t *testing.T, kindName string, client Client) {
	ctx := context.Background()

	want := "test-payload"
	src := sourceData{Payload: want}
	key, err := client.Put(ctx, &datastore.Key{Kind: kindName}, &src)
	if err != nil {
		t.Fatal(err)
	}

	var dst sourceData
	if err := client.Get(ctx, key, &dst); err != nil {
		t.Fatal(err)
	}

	if dst.Payload != want {
		t.Fatalf(`expected %q to equal %s`, dst.Payload, want)
	}

	if err := client.Delete(ctx, key); err != nil {
		t.Fatal(err)
	}

	if err := client.Get(ctx, key, &dst); err != datastore.ErrNoSuchEntity {
		t.Fatalf("expected ErrNoSuchEntity error: %v", err)
	}
}

// This test demonstrates how to use this package to create a simple fake for
// the datastore client.
func TestFake(t *testing.T) {
	client := newFakeClient()
	basicTests(t, "test-kind", client)
}

type fakeClient struct {
	Client
	m map[string]interface{}
}

func newFakeClient() Client {
	return &fakeClient{
		m: make(map[string]interface{}),
	}
}

func (c *fakeClient) Put(ctx context.Context, key *datastore.Key, src interface{}) (*datastore.Key, error) {
	if key.Name == "" && key.ID == 0 {
		key.ID = rand.Int63()
	}

	c.m[key.String()] = *(src.(*sourceData))
	return key, nil
}

func (c *fakeClient) Get(ctx context.Context, key *datastore.Key, dst interface{}) error {
	val, ok := c.m[key.String()]
	if !ok {
		return datastore.ErrNoSuchEntity
	}

	sd := dst.(*sourceData)
	sv := val.(sourceData)
	*sd = sv
	return nil
}

func (c *fakeClient) Delete(ctx context.Context, key *datastore.Key) error {
	delete(c.m, key.String())
	return nil
}
