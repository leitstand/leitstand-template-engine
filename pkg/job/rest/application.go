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
	"net/http"

	"github.com/leitstand/leitstand-template-engine/pkg/job"
	"github.com/leitstand/leitstand-template-engine/pkg/util"
	"github.com/rs/zerolog/log"
)

//Application is the port of CTRLD to get all the information of container images.
type Application struct {
	repository job.Repository
}

//NewApplication creates a new Application
func NewApplication(repository job.Repository) *Application {
	return &Application{
		repository: repository,
	}
}

//jobs
//@Summary "Jobs": list all jobs
//@Description Lists all jobs as list. Filtering is not possible.
//@Tags jobs
//@Accept  json
//@Produce  json
//@Success 200 {array} job.Job "list of jobs"
//@Failure 500 {object} util.Message
//@Router /template-engine/api/v1/jobs [get]
func (app *Application) jobs(w http.ResponseWriter, _ *http.Request) {
	jobs, err := app.repository.Jobs()
	if err != nil {
		log.Error().Err(err).Msg("error in getting Jobs")
		util.WriteMessage(w, http.StatusInternalServerError, "error in getting Jobs")
		return
	}
	jobList := make([]*job.Job, 0, len(jobs))
	for _, entity := range jobs {
		jobList = append(jobList, entity)
	}
	util.WriteAsJSON(w, http.StatusOK, jobList)
}

//job
//@Summary "Jobs": get a particular job
//@Description Returns a particular job.
//@Description **Characteristics:**
//@Description * Operation: **sync**
//@Tags jobs
//@Accept  json
//@Produce  json
//@Param id path string true "id of the job"
//@Success 200 {object} job.Result "Job Done"
//@Success 202 {object} job.Result "Operation is still pending"
//@Success 404 {object} util.Message "job not found"
//@Failure 500 {object} util.Message
//@Router /template-engine/api/v1/jobs/{id} [get]
func (app *Application) job(w http.ResponseWriter, req *http.Request) {
	id, ok := validateAndGetIDFromPath(w, req)
	if !ok {
		return
	}
	entity, err := app.repository.Job(id)
	if err != nil {
		log.Error().Err(err).Msg("error in getting Jobs")
		util.WriteMessage(w, http.StatusInternalServerError, "error in getting job")
		return
	}
	if entity == nil {
		util.WriteMessage(w, http.StatusNotFound, "job not found")
		return
	}
	if entity.State == job.StatusPending {
		util.WriteAsJSON(w, http.StatusAccepted, entity)
		return
	}
	util.WriteAsJSON(w, http.StatusOK, entity)
}
