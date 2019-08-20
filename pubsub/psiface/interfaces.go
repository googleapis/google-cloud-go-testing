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

	"cloud.google.com/go/pubsub"
)

type Client interface {
	CreateTopic(ctx context.Context, topicID string) (Topic, error)
	Topic(id string) Topic
	CreateSubscription(ctx context.Context, id string, topicID string) (Subscription, error)
	Subscription(id string) Subscription

	embedToIncludeNewMethods()
}

type Topic interface {
	String() string
	Publish(ctx context.Context, msg *pubsub.Message) PublishResult

	embedToIncludeNewMethods()
}

type Subscription interface {
	Exists(ctx context.Context) (bool, error)
	Receive(ctx context.Context, f func(context.Context, *pubsub.Message)) error
	Delete(ctx context.Context) error

	embedToIncludeNewMethods()
}

type PublishResult interface {
	Get(ctx context.Context) (serverID string, err error)

	embedToIncludeNewMethods()
}
