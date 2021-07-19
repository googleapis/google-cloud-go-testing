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
	"time"

	"cloud.google.com/go/pubsub"
)

// AdaptClient adapts a pubsub.Client so that it satisfies the Client
// interface.
func AdaptClient(c *pubsub.Client) Client {
	return client{c}
}

// AdaptMessage adapts a pubsub.Message so that it satisfies the Message
// interface.
func AdaptMessage(msg *pubsub.Message) Message {
	return message{msg}
}

type (
	client        struct{ *pubsub.Client }
	topic         struct{ *pubsub.Topic }
	subscription  struct{ *pubsub.Subscription }
	message       struct{ *pubsub.Message }
	publishResult struct{ *pubsub.PublishResult }
)

func (client) embedToIncludeNewMethods()        {}
func (topic) embedToIncludeNewMethods()         {}
func (subscription) embedToIncludeNewMethods()  {}
func (message) embedToIncludeNewMethods()       {}
func (publishResult) embedToIncludeNewMethods() {}

func (c client) CreateTopic(ctx context.Context, topicID string) (Topic, error) {
	t, err := c.Client.CreateTopic(ctx, topicID)
	if err != nil {
		return nil, err
	}
	return topic{t}, nil
}

func (c client) Topic(id string) Topic {
	return topic{c.Client.Topic(id)}
}

func (c client) CreateSubscription(ctx context.Context, id string, cfg SubscriptionConfig) (Subscription, error) {
	s, err := c.Client.CreateSubscription(ctx, id, cfg.toPS())
	if err != nil {
		return nil, err
	}
	return subscription{s}, nil
}

func (c client) Subscription(id string) Subscription {
	return subscription{c.Client.Subscription(id)}
}

func (t topic) String() string {
	return t.Topic.String()
}

func (t topic) Publish(ctx context.Context, msg Message) PublishResult {
	return publishResult{t.Topic.Publish(ctx, msg.(message).Message)}
}

func (s subscription) Exists(ctx context.Context) (bool, error) {
	return s.Subscription.Exists(ctx)
}

func (s subscription) Receive(ctx context.Context, f func(ctx context.Context, msg Message)) error {
	return s.Subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		f(ctx, AdaptMessage(msg))
	})
}

func (s subscription) Delete(ctx context.Context) error {
	return s.Subscription.Delete(ctx)
}

func (m message) ID() string {
	return m.Message.ID
}

func (m message) Data() []byte {
	return m.Message.Data
}

func (m message) Attributes() map[string]string {
	return m.Message.Attributes
}

func (m message) PublishTime() time.Time {
	return m.Message.PublishTime
}

func (r publishResult) Get(ctx context.Context) (serverID string, err error) {
	return r.PublishResult.Get(ctx)
}

func (cfg SubscriptionConfig) toPS() pubsub.SubscriptionConfig {
	return pubsub.SubscriptionConfig{
		Topic:               cfg.Topic.(topic).Topic,
		PushConfig:          cfg.PushConfig,
		AckDeadline:         cfg.AckDeadline,
		RetainAckedMessages: cfg.RetainAckedMessages,
		RetentionDuration:   cfg.RetentionDuration,
		Labels:              cfg.Labels,
	}
}
