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

	"cloud.google.com/go/bigquery"
)

func AdaptClient(c *bigquery.Client) Client {
	return client{c}
}

type (
	client          struct{ *bigquery.Client }
	copier          struct{ *bigquery.Copier }
	dataset         struct{ *bigquery.Dataset }
	datasetIterator struct{ *bigquery.DatasetIterator }
	extractor       struct{ *bigquery.Extractor }
	job             struct{ *bigquery.Job }
	jobIterator     struct{ *bigquery.JobIterator }
	loader          struct{ *bigquery.Loader }
	query           struct{ *bigquery.Query }
	rowIterator     struct{ *bigquery.RowIterator }
	table           struct{ *bigquery.Table }
	tableIterator   struct{ *bigquery.TableIterator }
	uploader        struct{ *bigquery.Uploader }
)

func (client) embedToIncludeNewMethods()          {}
func (copier) embedToIncludeNewMethods()          {}
func (dataset) embedToIncludeNewMethods()         {}
func (extractor) embedToIncludeNewMethods()       {}
func (job) embedToIncludeNewMethods()             {}
func (jobIterator) embedToIncludeNewMethods()     {}
func (loader) embedToIncludeNewMethods()          {}
func (query) embedToIncludeNewMethods()           {}
func (rowIterator) embedToIncludeNewMethods()     {}
func (table) embedToIncludeNewMethods()           {}
func (datasetIterator) embedToIncludeNewMethods() {}
func (tableIterator) embedToIncludeNewMethods()   {}
func (uploader) embedToIncludeNewMethods()        {}

func (c client) Location() string                     { return c.Client.Location }
func (c client) SetLocation(s string)                 { c.Client.Location = s }
func (c client) Close() error                         { return c.Client.Close() }
func (c client) Dataset(id string) Dataset            { return dataset{c.Client.Dataset(id)} }
func (c client) Jobs(ctx context.Context) JobIterator { return jobIterator{c.Client.Jobs(ctx)} }
func (c client) Query(s string) Query                 { return query{c.Client.Query(s)} }

func (c client) DatasetInProject(p, d string) Dataset {
	return dataset{c.Client.DatasetInProject(p, d)}
}

func (c client) Datasets(ctx context.Context) DatasetIterator {
	return datasetIterator{c.Client.Datasets(ctx)}
}

func (c client) DatasetsInProject(ctx context.Context, p string) DatasetIterator {
	return datasetIterator{c.Client.DatasetsInProject(ctx, p)}
}

func (c client) JobFromID(ctx context.Context, id string) (Job, error) {
	return adaptJob(c.Client.JobFromID(ctx, id))
}

func (c client) JobFromIDLocation(ctx context.Context, id, location string) (Job, error) {
	return adaptJob(c.Client.JobFromIDLocation(ctx, id, location))
}

func (c copier) JobIDConfig() *bigquery.JobIDConfig { return &c.Copier.JobIDConfig }

func (c copier) SetCopyConfig(cc CopyConfig) {
	c.Copier.CopyConfig = cc.CopyConfig
	for _, t := range cc.Srcs {
		c.Copier.CopyConfig.Srcs = append(c.Copier.CopyConfig.Srcs, t.(table).Table)
	}
	c.Copier.CopyConfig.Dst = cc.Dst.(table).Table
}

func (c copier) Run(ctx context.Context) (Job, error) {
	return adaptJob(c.Copier.Run(ctx))
}

func (d dataset) ProjectID() string                { return d.Dataset.ProjectID }
func (d dataset) DatasetID() string                { return d.Dataset.DatasetID }
func (d dataset) Delete(ctx context.Context) error { return d.Dataset.Delete(ctx) }
func (d dataset) DeleteWithContents(ctx context.Context) error {
	return d.Dataset.DeleteWithContents(ctx)
}
func (d dataset) Table(id string) Table { return table{d.Dataset.Table(id)} }
func (d dataset) Tables(ctx context.Context) TableIterator {
	return tableIterator{d.Dataset.Tables(ctx)}
}

func (d dataset) Create(ctx context.Context, dm *DatasetMetadata) error {
	return d.Dataset.Create(ctx, dm.toBQ())
}

func (d dataset) Metadata(ctx context.Context) (*DatasetMetadata, error) {
	m, err := d.Dataset.Metadata(ctx)
	return datasetMetadataFromBQ(m), err
}

func (d dataset) Update(ctx context.Context, dm DatasetMetadataToUpdate, etag string) (*DatasetMetadata, error) {
	m, err := d.Dataset.Update(ctx, dm.toBQ(), etag)
	if err != nil {
		return nil, err
	}
	return datasetMetadataFromBQ(m), nil
}

func (di datasetIterator) SetListHidden(b bool)  { di.DatasetIterator.ListHidden = b }
func (di datasetIterator) SetFilter(s string)    { di.DatasetIterator.Filter = s }
func (di datasetIterator) SetProjectID(s string) { di.DatasetIterator.ProjectID = s }

func (di datasetIterator) Next() (Dataset, error) {
	ds, err := di.DatasetIterator.Next()
	if err != nil {
		return nil, err
	}
	return dataset{ds}, nil
}

func (e extractor) JobIDConfig() *bigquery.JobIDConfig { return &e.Extractor.JobIDConfig }

func (e extractor) SetExtractConfig(c ExtractConfig) {
	e.Extractor.ExtractConfig = c.ExtractConfig
	e.Extractor.ExtractConfig.Src = c.Src.(table).Table
}

func (e extractor) Run(ctx context.Context) (Job, error) {
	return adaptJob(e.Extractor.Run(ctx))
}

func (j job) Status(ctx context.Context) (*bigquery.JobStatus, error) { return j.Job.Status(ctx) }
func (j job) LastStatus() *bigquery.JobStatus                         { return j.Job.LastStatus() }
func (j job) Cancel(ctx context.Context) error                        { return j.Job.Cancel(ctx) }
func (j job) Wait(ctx context.Context) (*bigquery.JobStatus, error)   { return j.Job.Wait(ctx) }

func (j job) Read(ctx context.Context) (RowIterator, error) {
	r, err := j.Job.Read(ctx)
	if err != nil {
		return nil, err
	}
	return rowIterator{r}, nil
}

func (j jobIterator) SetProjectID(s string)     { j.JobIterator.ProjectID = s }
func (j jobIterator) SetAllUsers(b bool)        { j.JobIterator.AllUsers = b }
func (j jobIterator) SetState(s bigquery.State) { j.JobIterator.State = s }

func (j jobIterator) Next() (Job, error) {
	return adaptJob(j.JobIterator.Next())
}

func (l loader) JobIDConfig() *bigquery.JobIDConfig { return &l.Loader.JobIDConfig }

func (l loader) SetLoadConfig(c LoadConfig) {
	l.Loader.LoadConfig = c.LoadConfig
	l.Loader.LoadConfig.Dst = c.Dst.(table).Table
}

func (l loader) Run(ctx context.Context) (Job, error) {
	return adaptJob(l.Loader.Run(ctx))
}

func (q query) JobIDConfig() *bigquery.JobIDConfig   { return &q.Query.JobIDConfig }
func (q query) Run(ctx context.Context) (Job, error) { return adaptJob(q.Query.Run(ctx)) }

func (q query) Read(ctx context.Context) (RowIterator, error) {
	r, err := q.Query.Read(ctx)
	if err != nil {
		return nil, err
	}
	return rowIterator{r}, nil
}

func (q query) SetQueryConfig(c QueryConfig) {
	q.Query.QueryConfig = c.QueryConfig
	if c.Dst != nil {
		q.Query.QueryConfig.Dst = c.Dst.(table).Table
	}
}

func (r rowIterator) SetStartIndex(i uint64)     { r.RowIterator.StartIndex = i }
func (r rowIterator) Schema() bigquery.Schema    { return r.RowIterator.Schema }
func (r rowIterator) TotalRows() uint64          { return r.RowIterator.TotalRows }
func (r rowIterator) Next(dst interface{}) error { return r.RowIterator.Next(dst) }

func (t table) ProjectID() string                    { return t.Table.ProjectID }
func (t table) DatasetID() string                    { return t.Table.DatasetID }
func (t table) TableID() string                      { return t.Table.TableID }
func (t table) FullyQualifiedName() string           { return t.Table.FullyQualifiedName() }
func (t table) Delete(ctx context.Context) error     { return t.Table.Delete(ctx) }
func (t table) Uploader() Uploader                   { return uploader{t.Table.Uploader()} }
func (t table) Read(ctx context.Context) RowIterator { return rowIterator{t.Table.Read(ctx)} }

func (t table) Create(ctx context.Context, tm *bigquery.TableMetadata) error {
	return t.Table.Create(ctx, tm)
}
func (t table) Metadata(ctx context.Context) (*bigquery.TableMetadata, error) {
	return t.Table.Metadata(ctx)
}

func (t table) Update(ctx context.Context, tm bigquery.TableMetadataToUpdate, etag string) (*bigquery.TableMetadata, error) {
	return t.Table.Update(ctx, tm, etag)
}

func (t table) CopierFrom(ts ...Table) Copier {
	var bts []*bigquery.Table
	for _, tb := range ts {
		bts = append(bts, tb.(table).Table)
	}
	c := t.Table.CopierFrom(bts...)
	return copier{c}
}

func (t table) ExtractorTo(dst *bigquery.GCSReference) Extractor {
	return extractor{t.Table.ExtractorTo(dst)}
}

func (t table) LoaderFrom(s bigquery.LoadSource) Loader {
	return loader{t.Table.LoaderFrom(s)}
}

func (ti tableIterator) Next() (Table, error) {
	t, err := ti.TableIterator.Next()
	if err != nil {
		return nil, err
	}
	return table{t}, nil
}

func (u uploader) SetSkipInvalidRows(b bool)                    { u.Uploader.SkipInvalidRows = b }
func (u uploader) SetIgnoreUnknownValues(b bool)                { u.Uploader.IgnoreUnknownValues = b }
func (u uploader) SetTableTemplateSuffix(s string)              { u.Uploader.TableTemplateSuffix = s }
func (u uploader) Put(ctx context.Context, i interface{}) error { return u.Uploader.Put(ctx, i) }

func adaptJob(j *bigquery.Job, err error) (Job, error) {
	if err != nil {
		return nil, err
	}
	return job{j}, nil
}

func (m *DatasetMetadata) toBQ() *bigquery.DatasetMetadata {
	m.DatasetMetadata.Access = accessEntriesToBQ(m.Access)
	return &m.DatasetMetadata
}

func datasetMetadataFromBQ(m *bigquery.DatasetMetadata) *DatasetMetadata {
	if m == nil {
		return nil
	}
	return &DatasetMetadata{
		DatasetMetadata: *m,
		Access:          accessEntriesFromBQ(m.Access),
	}
}

func (u DatasetMetadataToUpdate) toBQ() bigquery.DatasetMetadataToUpdate {
	u.DatasetMetadataToUpdate.Access = accessEntriesToBQ(u.Access)
	return u.DatasetMetadataToUpdate
}

func (e *AccessEntry) toBQ() *bigquery.AccessEntry {
	if e.View != nil {
		e.AccessEntry.View = e.View.(table).Table
	}
	return &e.AccessEntry
}

func accessEntryFromBQ(e *bigquery.AccessEntry) *AccessEntry {
	if e == nil {
		return nil
	}
	r := &AccessEntry{AccessEntry: *e}
	if e.View != nil {
		r.View = table{e.View}
	}
	return r
}

func accessEntriesToBQ(a []*AccessEntry) []*bigquery.AccessEntry {
	var r []*bigquery.AccessEntry
	for _, e := range a {
		r = append(r, e.toBQ())
	}
	return r
}

func accessEntriesFromBQ(a []*bigquery.AccessEntry) []*AccessEntry {
	var r []*AccessEntry
	for _, e := range a {
		r = append(r, accessEntryFromBQ(e))
	}
	return r
}
