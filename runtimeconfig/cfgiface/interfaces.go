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

package cfgiface

import (
	"google.golang.org/api/googleapi"
	runtimeconfig "google.golang.org/api/runtimeconfig/v1beta1"
)

type Service interface {
	Projects() ProjectsService

	embedToIncludeNewMethods()
}

type ProjectsService interface {
	Configs() ProjectsConfigsService

	embedToIncludeNewMethods()
}

type ProjectsConfigsService interface {
	Variables() ProjectsConfigsVariablesService

	embedToIncludeNewMethods()
}

type ProjectsConfigsVariablesService interface {
	Create(parent string, variable *runtimeconfig.Variable) ProjectsConfigsVariablesCreateCall
	Delete(name string) ProjectsConfigsVariablesDeleteCall
	List(parent string) ProjectsConfigsVariablesListCall
	Update(name string, variable *runtimeconfig.Variable) ProjectsConfigsVariablesUpdateCall
	Watch(name string, watchvariablerequest *runtimeconfig.WatchVariableRequest) ProjectsConfigsVariablesWatchCall

	embedToIncludeNewMethods()
}

type ProjectsConfigsVariablesCreateCall interface {
	Do(opts ...googleapi.CallOption) (*runtimeconfig.Variable, error)

	embedToIncludeNewMethods()
}

type ProjectsConfigsVariablesDeleteCall interface {
	Recursive(recursive bool) ProjectsConfigsVariablesDeleteCall
	Do(opts ...googleapi.CallOption) (*runtimeconfig.Empty, error)

	embedToIncludeNewMethods()
}

type ProjectsConfigsVariablesListCall interface {
	Do(opts ...googleapi.CallOption) (*runtimeconfig.ListVariablesResponse, error)

	embedToIncludeNewMethods()
}

type ProjectsConfigsVariablesUpdateCall interface {
	Do(opts ...googleapi.CallOption) (*runtimeconfig.Variable, error)

	embedToIncludeNewMethods()
}

type ProjectsConfigsVariablesWatchCall interface {
	Do(opts ...googleapi.CallOption) (*runtimeconfig.Variable, error)

	embedToIncludeNewMethods()
}
