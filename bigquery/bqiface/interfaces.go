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
	"google.golang.org/api/iterator"
)

type Client interface {
	Location() string
	SetLocation(string)
	Close() error
	Dataset(string) Dataset
	DatasetInProject(string, string) Dataset
	Datasets(context.Context) DatasetIterator
	DatasetsInProject(context.Context, string) DatasetIterator
	Query(string) Query
	JobFromID(context.Context, string) (Job, error)
	JobFromIDLocation(context.Context, string, string) (Job, error)
	Jobs(context.Context) JobIterator

	embedToIncludeNewMethods()
}

type Copier interface {
	JobIDConfig() *bigquery.JobIDConfig
	SetCopyConfig(CopyConfig)
	Run(context.Context) (Job, error)

	embedToIncludeNewMethods()
}

type Dataset interface {
	ProjectID() string
	DatasetID() string
	Create(context.Context, *DatasetMetadata) error
	Delete(context.Context) error
	DeleteWithContents(context.Context) error
	Metadata(context.Context) (*DatasetMetadata, error)
	Update(context.Context, DatasetMetadataToUpdate, string) (*DatasetMetadata, error)
	Table(string) Table
	Tables(context.Context) TableIterator

	embedToIncludeNewMethods()
}

type DatasetIterator interface {
	SetListHidden(bool)
	SetFilter(string)
	SetProjectID(string)
	Next() (Dataset, error)
	PageInfo() *iterator.PageInfo

	embedToIncludeNewMethods()
}

type Extractor interface {
	JobIDConfig() *bigquery.JobIDConfig
	SetExtractConfig(ExtractConfig)
	Run(context.Context) (Job, error)

	embedToIncludeNewMethods()
}

type Loader interface {
	JobIDConfig() *bigquery.JobIDConfig
	SetLoadConfig(LoadConfig)
	Run(context.Context) (Job, error)

	embedToIncludeNewMethods()
}

type Job interface {
	ID() string
	Location() string
	Config() (bigquery.JobConfig, error)
	Status(context.Context) (*bigquery.JobStatus, error)
	LastStatus() *bigquery.JobStatus
	Cancel(context.Context) error
	Wait(context.Context) (*bigquery.JobStatus, error)
	Read(context.Context) (RowIterator, error)

	embedToIncludeNewMethods()
}

type JobIterator interface {
	SetProjectID(string)
	SetAllUsers(bool)
	SetState(bigquery.State)
	Next() (Job, error)
	PageInfo() *iterator.PageInfo

	embedToIncludeNewMethods()
}

type Query interface {
	JobIDConfig() *bigquery.JobIDConfig
	SetQueryConfig(QueryConfig)
	Run(context.Context) (Job, error)
	Read(context.Context) (RowIterator, error)

	embedToIncludeNewMethods()
}

type RowIterator interface {
	SetStartIndex(uint64)
	Schema() bigquery.Schema
	TotalRows() uint64
	Next(interface{}) error
	PageInfo() *iterator.PageInfo

	embedToIncludeNewMethods()
}

type Table interface {
	CopierFrom(...Table) Copier
	Create(context.Context, *bigquery.TableMetadata) error
	DatasetID() string
	Delete(context.Context) error
	ExtractorTo(dst *bigquery.GCSReference) Extractor
	FullyQualifiedName() string
	LoaderFrom(bigquery.LoadSource) Loader
	Metadata(context.Context) (*bigquery.TableMetadata, error)
	ProjectID() string
	Read(ctx context.Context) RowIterator
	TableID() string
	Update(context.Context, bigquery.TableMetadataToUpdate, string) (*bigquery.TableMetadata, error)
	Uploader() Uploader

	embedToIncludeNewMethods()
}

type TableIterator interface {
	Next() (Table, error)
	PageInfo() *iterator.PageInfo

	embedToIncludeNewMethods()
}

type Uploader interface {
	SetSkipInvalidRows(bool)
	SetIgnoreUnknownValues(bool)
	SetTableTemplateSuffix(string)
	Put(context.Context, interface{}) error

	embedToIncludeNewMethods()
}
