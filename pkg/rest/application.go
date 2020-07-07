/*
 * Copyright 2020 RtBrick Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 * use this file except in compliance with the License.  You may obtain a copy
 * of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  See the
 * License for the specific language governing permissions and limitations under
 * the License.
 */

package rest

import (
	"github.com/leitstand/leitstand-template-engine/pkg/configen"
	"github.com/leitstand/leitstand-template-engine/pkg/job"
)

// Application is the configen application
type Application struct {
	repository    *configen.Repository
	jobRepository job.Repository
}

// NewApplication creates a new Application
func NewApplication(repository *configen.Repository, jobRepository job.Repository) *Application {
	return &Application{
		repository:    repository,
		jobRepository: jobRepository,
	}
}
