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

// AdaptClient adapts a datastore.Client so that it satisfies the Client
// interface.
func AdaptClient(c *datastore.Client) Client {
	return client{c}
}

type (
	client      struct{ *datastore.Client }
	transaction struct{ *datastore.Transaction }
	iterator    struct{ *datastore.Iterator }
	commit      struct{ *datastore.Commit }
)

func (client) embedToIncludeNewMethods()      {}
func (transaction) embedToIncludeNewMethods() {}
func (iterator) embedToIncludeNewMethods()    {}
func (commit) embedToIncludeNewMethods()      {}

func (c client) Close() error {
	return c.Client.Close()
}

func (c client) AllocateIDs(ctx context.Context, keys []*datastore.Key) ([]*datastore.Key, error) {
	return c.Client.AllocateIDs(ctx, keys)
}

func (c client) Count(ctx context.Context, q *datastore.Query) (int, error) {
	return c.Client.Count(ctx, q)
}

func (c client) Delete(ctx context.Context, key *datastore.Key) error {
	return c.Client.Delete(ctx, key)
}

func (c client) DeleteMulti(ctx context.Context, keys []*datastore.Key) error {
	return c.Client.DeleteMulti(ctx, keys)
}

func (c client) Get(ctx context.Context, key *datastore.Key, dst interface{}) error {
	return c.Client.Get(ctx, key, dst)
}

func (c client) GetAll(ctx context.Context, q *datastore.Query, dst interface{}) ([]*datastore.Key, error) {
	return c.Client.GetAll(ctx, q, dst)
}

func (c client) GetMulti(ctx context.Context, keys []*datastore.Key, dst interface{}) error {
	return c.Client.GetMulti(ctx, keys, dst)
}

func (c client) Mutate(ctx context.Context, muts ...*datastore.Mutation) ([]*datastore.Key, error) {
	return c.Client.Mutate(ctx, muts...)
}

func (c client) NewTransaction(ctx context.Context, opts ...datastore.TransactionOption) (Transaction, error) {
	t, err := c.Client.NewTransaction(ctx, opts...)
	return transaction{t}, err
}

func (c client) Put(ctx context.Context, key *datastore.Key, src interface{}) (*datastore.Key, error) {
	return c.Client.Put(ctx, key, src)
}

func (c client) PutMulti(ctx context.Context, keys []*datastore.Key, src interface{}) ([]*datastore.Key, error) {
	return c.Client.PutMulti(ctx, keys, src)
}

func (c client) Run(ctx context.Context, q *datastore.Query) Iterator {
	return iterator{c.Client.Run(ctx, q)}
}

func (c client) RunInTransaction(ctx context.Context, f func(tx Transaction) error, opts ...datastore.TransactionOption) (Commit, error) {
	cmt, err := c.Client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		return f(transaction{tx})
	}, opts...)
	return commit{cmt}, err
}

func (t transaction) Commit() (Commit, error) {
	c, err := t.Transaction.Commit()
	return commit{c}, err
}

func (t transaction) Delete(key *datastore.Key) error {
	return t.Transaction.Delete(key)
}

func (t transaction) DeleteMulti(keys []*datastore.Key) error {
	return t.Transaction.DeleteMulti(keys)
}

func (t transaction) Get(key *datastore.Key, dst interface{}) error {
	return t.Transaction.Get(key, dst)
}

func (t transaction) GetMulti(keys []*datastore.Key, dst interface{}) error {
	return t.Transaction.GetMulti(keys, dst)
}

func (t transaction) Mutate(muts ...*datastore.Mutation) ([]*datastore.PendingKey, error) {
	return t.Transaction.Mutate(muts...)
}

func (t transaction) Put(key *datastore.Key, src interface{}) (*datastore.PendingKey, error) {
	return t.Transaction.Put(key, src)
}

func (t transaction) PutMulti(keys []*datastore.Key, src interface{}) ([]*datastore.PendingKey, error) {
	return t.Transaction.PutMulti(keys, src)
}

func (t transaction) Rollback() error {
	return t.Transaction.Rollback()
}

func (c commit) Key(p *datastore.PendingKey) *datastore.Key {
	return c.Commit.Key(p)
}
