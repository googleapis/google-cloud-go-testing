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

	"cloud.google.com/go/datastore"
)

// Client is the interface that wraps a datastore.Client.
type Client interface {
	Close() error
	AllocateIDs(ctx context.Context, keys []*datastore.Key) ([]*datastore.Key, error)
	Count(ctx context.Context, q *datastore.Query) (n int, err error)
	Delete(ctx context.Context, key *datastore.Key) error
	DeleteMulti(ctx context.Context, keys []*datastore.Key) (err error)
	Get(ctx context.Context, key *datastore.Key, dst interface{}) (err error)
	GetAll(ctx context.Context, q *datastore.Query, dst interface{}) (keys []*datastore.Key, err error)
	GetMulti(ctx context.Context, keys []*datastore.Key, dst interface{}) (err error)
	Mutate(ctx context.Context, muts ...*datastore.Mutation) (ret []*datastore.Key, err error)
	NewTransaction(ctx context.Context, opts ...datastore.TransactionOption) (t Transaction, err error)
	Put(ctx context.Context, key *datastore.Key, src interface{}) (*datastore.Key, error)
	PutMulti(ctx context.Context, keys []*datastore.Key, src interface{}) (ret []*datastore.Key, err error)
	Run(ctx context.Context, q *datastore.Query) Iterator
	RunInTransaction(ctx context.Context, f func(tx Transaction) error, opts ...datastore.TransactionOption) (cmt Commit, err error)

	embedToIncludeNewMethods()
}

// Transaction is the interface that wraps a datastore.Transaction.
type Transaction interface {
	Commit() (c Commit, err error)
	Delete(key *datastore.Key) error
	DeleteMulti(keys []*datastore.Key) (err error)
	Get(key *datastore.Key, dst interface{}) (err error)
	GetMulti(keys []*datastore.Key, dst interface{}) (err error)
	Mutate(muts ...*datastore.Mutation) ([]*datastore.PendingKey, error)
	Put(key *datastore.Key, src interface{}) (*datastore.PendingKey, error)
	PutMulti(keys []*datastore.Key, src interface{}) (ret []*datastore.PendingKey, err error)
	Rollback() (err error)

	embedToIncludeNewMethods()
}

// Iterator is the interface that wraps a datastore.Iterator.
type Iterator interface {
	Cursor() (c datastore.Cursor, err error)
	Next(dst interface{}) (k *datastore.Key, err error)

	embedToIncludeNewMethods()
}

// Commit is the interface that wraps a datastore.Commit.
type Commit interface {
	Key(p *datastore.PendingKey) *datastore.Key

	embedToIncludeNewMethods()
}
