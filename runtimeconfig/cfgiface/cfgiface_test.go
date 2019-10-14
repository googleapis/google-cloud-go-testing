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
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"google.golang.org/api/googleapi"
	runtimeconfig "google.golang.org/api/runtimeconfig/v1beta1"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("integration tests skipped in short mode")
	}

	configName := os.Getenv("CFGIFACE_CONFIG")
	if configName == "" {
		t.Skip("missing CFGIFACE_CONFIG environment variable")
	}

	ctx := context.Background()

	s, err := runtimeconfig.NewService(ctx)
	if err != nil {
		t.Fatal(err)
	}

	service := AdaptService(s)

	basicTests(t, configName, service)
}

func basicTests(t *testing.T, configName string, service Service) {
	op := service.Projects().Configs().Variables()

	variable := &runtimeconfig.Variable{
		Name: fmt.Sprintf("%v/variables/cfgiface_test/%d", configName, time.Now().UnixNano()),
		Text: "hello, cfgiface",
	}

	_, err := op.Delete(path.Dir(variable.Name)).Recursive(true).Do()
	if err != nil {
		t.Fatal(err)
	}

	if variable, err = op.Create(configName, variable).Do(); err != nil {
		t.Fatal(err)
	}

	response, err := op.List(configName).Do()
	if err != nil {
		t.Fatal(err)
	}

	vars := make(map[string]*runtimeconfig.Variable, len(response.Variables))
	for _, v := range response.Variables {
		vars[v.Name] = v
	}

	if _, ok := vars[variable.Name]; !ok {
		t.Fatal("variable not returned in list")
	}

	changes := make(chan *runtimeconfig.Variable, 1)
	errs := make(chan error, 1)

	go func() {
		if change, err := op.Watch(variable.Name, &runtimeconfig.WatchVariableRequest{}).Do(); err != nil {
			errs <- err
		} else {
			changes <- change
		}
	}()

	variable.Text = "hello again, cfgiface"

	if variable, err = op.Update(variable.Name, variable).Do(); err != nil {
		t.Fatal(err)
	}

	select {
	case err := <-errs:
		t.Fatal(err)
	case change := <-changes:
		variable = change
	}

	if got, want := variable.State, "UPDATED"; got != want {
		t.Fatalf("state: got %v, want %v", got, want)
	}
}

// This test demonstrates how to use this package to create a simple fake for
// the runtimeconfig service.
func TestFake(t *testing.T) {
	basicTests(t, "my-config", newFakeService())
}

type fakeService struct {
	Service

	projects *fakeProjectsService
}

func newFakeService() Service {
	return &fakeService{}
}

func (s *fakeService) Projects() ProjectsService {
	if s.projects == nil {
		s.projects = &fakeProjectsService{
			configs: &fakeProjectsConfigsService{
				variables: &fakeProjectsConfigsVariablesService{
					vars:    make(map[string]*runtimeconfig.Variable),
					changes: make(chan *runtimeconfig.Variable),
				},
			},
		}
	}

	return s.projects
}

type fakeProjectsService struct {
	ProjectsService

	configs *fakeProjectsConfigsService
}

func (p *fakeProjectsService) Configs() ProjectsConfigsService {
	return p.configs
}

type fakeProjectsConfigsService struct {
	ProjectsConfigsService

	variables *fakeProjectsConfigsVariablesService
}

func (c *fakeProjectsConfigsService) Variables() ProjectsConfigsVariablesService {
	return c.variables
}

type fakeProjectsConfigsVariablesService struct {
	ProjectsConfigsVariablesService

	vars    map[string]*runtimeconfig.Variable
	changes chan *runtimeconfig.Variable
}

func (v *fakeProjectsConfigsVariablesService) Create(parent string, variable *runtimeconfig.Variable) ProjectsConfigsVariablesCreateCall {
	return &fakeProjectsConfigsVariablesCreateCall{parent: parent, variable: variable, svc: v}
}

type fakeProjectsConfigsVariablesCreateCall struct {
	ProjectsConfigsVariablesCreateCall

	parent   string
	variable *runtimeconfig.Variable

	svc *fakeProjectsConfigsVariablesService
}

func (c *fakeProjectsConfigsVariablesCreateCall) Do(opts ...googleapi.CallOption) (*runtimeconfig.Variable, error) {
	c.svc.vars[c.variable.Name] = c.variable
	return c.variable, nil
}

func (v *fakeProjectsConfigsVariablesService) Delete(name string) ProjectsConfigsVariablesDeleteCall {
	return &fakeProjectsConfigsVariablesDeleteCall{name: name, svc: v}
}

type fakeProjectsConfigsVariablesDeleteCall struct {
	ProjectsConfigsVariablesDeleteCall

	name      string
	recursive bool

	svc *fakeProjectsConfigsVariablesService
}

func (d *fakeProjectsConfigsVariablesDeleteCall) Recursive(recursive bool) ProjectsConfigsVariablesDeleteCall {
	d.recursive = recursive
	return d
}

func (d *fakeProjectsConfigsVariablesDeleteCall) Do(opts ...googleapi.CallOption) (*runtimeconfig.Empty, error) {
	switch {
	case d.recursive:
		for k := range d.svc.vars {
			if !strings.HasPrefix(k, d.name) {
				continue
			}

			delete(d.svc.vars, k)
		}
	default:
		delete(d.svc.vars, d.name)
	}

	return &runtimeconfig.Empty{}, nil
}

func (v *fakeProjectsConfigsVariablesService) List(parent string) ProjectsConfigsVariablesListCall {
	return &fakeProjectsConfigsVariablesListCall{svc: v}
}

type fakeProjectsConfigsVariablesListCall struct {
	ProjectsConfigsVariablesListCall

	svc *fakeProjectsConfigsVariablesService
}

func (l *fakeProjectsConfigsVariablesListCall) Do(opts ...googleapi.CallOption) (*runtimeconfig.ListVariablesResponse, error) {
	var r runtimeconfig.ListVariablesResponse
	for _, v := range l.svc.vars {
		r.Variables = append(r.Variables, v)
	}

	return &r, nil
}

func (v *fakeProjectsConfigsVariablesService) Update(name string, variable *runtimeconfig.Variable) ProjectsConfigsVariablesUpdateCall {
	return &fakeProjectsConfigsVariablesUpdateCall{name: name, variable: variable, svc: v}
}

type fakeProjectsConfigsVariablesUpdateCall struct {
	ProjectsConfigsVariablesUpdateCall

	name     string
	variable *runtimeconfig.Variable

	svc *fakeProjectsConfigsVariablesService
}

func (u *fakeProjectsConfigsVariablesUpdateCall) Do(opts ...googleapi.CallOption) (*runtimeconfig.Variable, error) {
	u.svc.vars[u.variable.Name] = u.variable
	u.variable.State = "UPDATED"
	u.svc.changes <- u.variable

	return u.variable, nil
}

func (v *fakeProjectsConfigsVariablesService) Watch(name string, watchvariablerequest *runtimeconfig.WatchVariableRequest) ProjectsConfigsVariablesWatchCall {
	return &fakeProjectsConfigsVariablesWatchCall{name: name, svc: v}
}

type fakeProjectsConfigsVariablesWatchCall struct {
	ProjectsConfigsVariablesWatchCall

	name string

	svc *fakeProjectsConfigsVariablesService
}

func (w *fakeProjectsConfigsVariablesWatchCall) Do(opts ...googleapi.CallOption) (*runtimeconfig.Variable, error) {
	for v := range w.svc.changes {
		if v.Name == w.name {
			return v, nil
		}
	}

	return nil, nil
}
