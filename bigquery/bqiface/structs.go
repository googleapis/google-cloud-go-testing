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
	"cloud.google.com/go/bigquery"
)

type AccessEntry struct {
	bigquery.AccessEntry
	View Table // shadows bigquery.AccessEntry's field
}

type CopyConfig struct {
	bigquery.CopyConfig
	Srcs []Table // shadows bigquery.CopyConfig's field
	Dst  Table   // shadows bigquery.CopyConfig's field
}

type DatasetMetadata struct {
	bigquery.DatasetMetadata
	Access []*AccessEntry // shadows bigquery.DatasetMetadata's field
}

type DatasetMetadataToUpdate struct {
	bigquery.DatasetMetadataToUpdate
	Access []*AccessEntry // shadows bigquery.DatasetMetadataToUpdate's field
}

type ExtractConfig struct {
	bigquery.ExtractConfig
	Src Table // shadows bigquery.ExtractConfig's field
}

type LoadConfig struct {
	bigquery.LoadConfig
	Dst Table // shadows bigquery.LoadConfig's field
}

type QueryConfig struct {
	bigquery.QueryConfig
	Dst Table // shaodws bigquery.QueryConfig's field
}
