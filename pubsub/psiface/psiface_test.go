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

package psiface

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("integration tests skipped in short mode")
	}

	topicID := os.Getenv("PSIFACE_TOPIC")
	if topicID == "" {
		t.Skip("missing PSIFACE_TOPIC environment variable")
	}
	projID, topicName, err := parseTopic(topicID)
	if err != nil {
		t.Fatal(err)
	}

	subscriptionName := fmt.Sprintf("psiface_test_%d", time.Now().UnixNano())
	ctx := context.Background()
	c, err := pubsub.NewClient(ctx, projID)
	if err != nil {
		t.Fatal(err)
	}
	client := AdaptClient(c)
	basicTests(t, topicName, subscriptionName, client)
}

func basicTests(t *testing.T, topicName string, subscriptionName string, client Client) {
	ctx := context.Background()
	topic := client.Topic(topicName)

	sub, err := client.CreateSubscription(ctx, subscriptionName, SubscriptionConfig{Topic: topic})
	if err != nil {
		t.Fatal(err)
	}

	const contents = "hello, psiface"
	ctx, cancel := context.WithCancel(ctx)
	errs := make(chan error, 50)
	go func() {
		err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			got, want := string(msg.Data), contents
			msg.Ack()
			if got == want {
				errs <- nil
				cancel()
			}
		})
		if err != nil {
			errs <- err
		}
	}()

	_, err = topic.Publish(ctx, &pubsub.Message{Data: []byte(contents)}).Get(ctx)
	if err != nil {
		t.Fatal(err)
	}
	err = <-errs
	if err != nil {
		t.Fatal(err)
	}

	ctx = context.Background()
	err = sub.Delete(ctx)
	if err != nil {
		t.Errorf("deleting: %v", err)
	}
}

func parseTopic(topicID string) (project, topic string, err error) {
	segs := strings.Split(topicID, "/")
	if len(segs) != 4 || segs[0] != "projects" || segs[2] != "topics" {
		return "", "", errors.New("invalid topic id")
	}
	return segs[1], segs[3], nil
}

// This test demonstrates how to use this package to create a simple fake for
// the pubsub client.
func TestFake(t *testing.T) {
	ctx := context.Background()
	client := newFakeClient()
	if _, err := client.CreateTopic(ctx, "my-topic"); err != nil {
		t.Fatal(err)
	}
	basicTests(t, "my-topic", "my-subscription", client)
}

type fakeClient struct {
	Client
	topics sync.Map
	subs   sync.Map
}

func newFakeClient() Client {
	return &fakeClient{}
}

func (c *fakeClient) CreateTopic(_ context.Context, topicID string) (Topic, error) {
	if _, ok := c.topics.Load(topicID); ok {
		return nil, fmt.Errorf("topic %q already exists", topicID)
	}
	t := &fakeTopic{c: c, name: topicID}
	c.topics.Store(topicID, t)
	return t, nil
}

func (c *fakeClient) Topic(id string) Topic {
	t, ok := c.topics.Load(id)
	if !ok {
		return &fakeTopic{c: c, name: id}
	}
	return t.(Topic)
}

func (c *fakeClient) CreateSubscription(ctx context.Context, id string, cfg SubscriptionConfig) (Subscription, error) {
	if _, ok := c.subs.Load(id); ok {
		return nil, fmt.Errorf("subscription %q already exists", id)
	}
	s := &fakeSubscription{
		c:       c,
		name:    id,
		topicID: cfg.Topic.String(),
		msgs:    make(chan *pubsub.Message, 50),
	}
	c.subs.Store(id, s)
	t := cfg.Topic.(*fakeTopic)
	t.subs = append(t.subs, s)
	return s, nil
}

func (c *fakeClient) Subscription(id string) Subscription {
	t, ok := c.subs.Load(id)
	if !ok {
		return &fakeSubscription{c: c, name: id}
	}
	return t.(Subscription)
}

type fakeTopic struct {
	Topic
	c    *fakeClient
	name string
	subs []*fakeSubscription
}

func (t *fakeTopic) String() string {
	return t.name
}

func (t *fakeTopic) Publish(ctx context.Context, msg *pubsub.Message) PublishResult {
	for _, sub := range t.subs {
		if sub.topicID == t.name {
			sub.msgs <- msg
		}
	}
	return &fakePublishResult{}
}

type fakeSubscription struct {
	Subscription
	c       *fakeClient
	name    string
	topicID string
	msgs    chan *pubsub.Message
}

func (s *fakeSubscription) Exists(_ context.Context) (bool, error) {
	_, ok := s.c.subs.Load(s.name)
	return ok, nil
}

func (s *fakeSubscription) Receive(ctx context.Context, f func(context.Context, *pubsub.Message)) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-s.msgs:
			if !ok {
				return nil
			}
			f(ctx, msg)
		}
	}
}

func (s *fakeSubscription) Delete(_ context.Context) error {
	s.c.subs.Delete(s.name)
	return nil
}

type fakePublishResult struct {
	PublishResult
}

func (r *fakePublishResult) Get(_ context.Context) (serverID string, err error) {
	return "", nil
}
