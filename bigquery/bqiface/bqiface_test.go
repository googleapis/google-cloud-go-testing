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

package bqiface

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("integration tests skipped in short mode")
	}
	projectID := os.Getenv("BQIFACE_PROJECT")
	if projectID == "" {
		t.Skip("missing BQIFACE_PROJECT environment variable")
	}

	ctx := context.Background()
	c, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		t.Fatal(err)
	}
	client := AdaptClient(c)
	defer client.Close()

	ds := client.Dataset(fmt.Sprintf("bqiface_%d", time.Now().Unix()))
	var wantMD DatasetMetadata
	wantMD.DefaultTableExpiration = time.Hour
	var ae AccessEntry
	ae.Role = bigquery.OwnerRole
	ae.EntityType = bigquery.SpecialGroupEntity
	ae.Entity = "projectOwners"
	wantMD.Access = []*AccessEntry{&ae}
	err = ds.Create(ctx, &wantMD)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := ds.Delete(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	gotMD, err := ds.Metadata(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := gotMD.DefaultTableExpiration, wantMD.DefaultTableExpiration; got != want {
		t.Errorf("DefaultTableExpiration: got %s, want %s", got, want)
	}
	if got, want := len(gotMD.Access), 1; got != want {
		t.Fatalf("got %d access entries, want %d", got, want)
	}
	if got, want := *gotMD.Access[0], ae; got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}

	table := ds.Table("t")
	schema := bigquery.Schema{
		{Name: "name", Type: bigquery.StringFieldType},
		{Name: "score", Type: bigquery.IntegerFieldType},
	}
	if err := table.Create(ctx, &bigquery.TableMetadata{Schema: schema}); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := table.Delete(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	upl := table.Uploader()
	var saverRows []*bigquery.ValuesSaver

	for i, name := range []string{"a", "b", "c"} {
		saverRows = append(saverRows, &bigquery.ValuesSaver{
			Schema:   schema,
			InsertID: name,
			Row:      []bigquery.Value{name, i},
		})
	}
	if err := upl.Put(ctx, saverRows); err != nil {
		t.Fatal(putError(err))
	}
	count := 0
	for {
		it := table.Read(ctx)
		count, err = countRows(it)
		if err != nil {
			t.Fatal(err)
		}
		if count > 0 {
			break
		}
		// Wait for rows to appear; it may take a few seconds.
		time.Sleep(1 * time.Second)
	}
	if got, want := count, len(saverRows); got != want {
		t.Errorf("got %d rows, want %d", got, want)
	}

	q := client.Query(fmt.Sprintf("SELECT * FROM %s.%s", table.DatasetID(), table.TableID()))
	it, err := q.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}
	count, err = countRows(it)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := count, len(saverRows); got != want {
		t.Errorf("got %d rows, want %d", got, want)
	}
}

func countRows(it RowIterator) (int, error) {
	n := 0
	for {
		var v []bigquery.Value
		err := it.Next(&v)
		if err == iterator.Done {
			return n, nil
		}
		if err != nil {
			return 0, err
		}
		n++
	}
}

func putError(err error) string {
	pme, ok := err.(bigquery.PutMultiError)
	if !ok {
		return err.Error()
	}
	var msgs []string
	for _, err := range pme {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "\n")
}
