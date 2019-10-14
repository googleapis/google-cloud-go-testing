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
	runtimeconfig "google.golang.org/api/runtimeconfig/v1beta1"
)

// AdaptClient adapts a runtimeconfig.Service so that it satisfies the Service
// interface.
func AdaptService(s *runtimeconfig.Service) Service {
	return service{s}
}

type (
	service                struct{ *runtimeconfig.Service }
	projectsService        struct{ *runtimeconfig.ProjectsService }
	projectsConfigsService struct {
		*runtimeconfig.ProjectsConfigsService
	}
	projectsConfigsVariablesService struct {
		*runtimeconfig.ProjectsConfigsVariablesService
	}
	projectsConfigsVariablesCreateCall struct {
		*runtimeconfig.ProjectsConfigsVariablesCreateCall
	}
	projectsConfigsVariablesDeleteCall struct {
		*runtimeconfig.ProjectsConfigsVariablesDeleteCall
	}
	projectsConfigsVariablesListCall struct {
		*runtimeconfig.ProjectsConfigsVariablesListCall
	}
	projectsConfigsVariablesUpdateCall struct {
		*runtimeconfig.ProjectsConfigsVariablesUpdateCall
	}
	projectsConfigsVariablesWatchCall struct {
		*runtimeconfig.ProjectsConfigsVariablesWatchCall
	}
)

func (service) embedToIncludeNewMethods()                            {}
func (projectsService) embedToIncludeNewMethods()                    {}
func (projectsConfigsService) embedToIncludeNewMethods()             {}
func (projectsConfigsVariablesService) embedToIncludeNewMethods()    {}
func (projectsConfigsVariablesCreateCall) embedToIncludeNewMethods() {}
func (projectsConfigsVariablesDeleteCall) embedToIncludeNewMethods() {}
func (projectsConfigsVariablesListCall) embedToIncludeNewMethods()   {}
func (projectsConfigsVariablesUpdateCall) embedToIncludeNewMethods() {}
func (projectsConfigsVariablesWatchCall) embedToIncludeNewMethods()  {}

func (s service) Projects() ProjectsService {
	return projectsService{s.Service.Projects}
}

func (p projectsService) Configs() ProjectsConfigsService {
	return projectsConfigsService{p.ProjectsService.Configs}
}

func (c projectsConfigsService) Variables() ProjectsConfigsVariablesService {
	return projectsConfigsVariablesService{c.ProjectsConfigsService.Variables}
}

func (v projectsConfigsVariablesService) Create(parent string, variable *runtimeconfig.Variable) ProjectsConfigsVariablesCreateCall {
	return projectsConfigsVariablesCreateCall{v.ProjectsConfigsVariablesService.Create(parent, variable)}
}

func (v projectsConfigsVariablesService) Delete(name string) ProjectsConfigsVariablesDeleteCall {
	return projectsConfigsVariablesDeleteCall{v.ProjectsConfigsVariablesService.Delete(name)}
}

func (d projectsConfigsVariablesDeleteCall) Recursive(recursive bool) ProjectsConfigsVariablesDeleteCall {
	return projectsConfigsVariablesDeleteCall{d.ProjectsConfigsVariablesDeleteCall.Recursive(recursive)}
}

func (v projectsConfigsVariablesService) List(parent string) ProjectsConfigsVariablesListCall {
	return projectsConfigsVariablesListCall{v.ProjectsConfigsVariablesService.List(parent)}
}

func (v projectsConfigsVariablesService) Update(name string, variable *runtimeconfig.Variable) ProjectsConfigsVariablesUpdateCall {
	return projectsConfigsVariablesUpdateCall{v.ProjectsConfigsVariablesService.Update(name, variable)}
}

func (v projectsConfigsVariablesService) Watch(name string, watchvariablerequest *runtimeconfig.WatchVariableRequest) ProjectsConfigsVariablesWatchCall {
	return projectsConfigsVariablesWatchCall{v.ProjectsConfigsVariablesService.Watch(name, watchvariablerequest)}
}
