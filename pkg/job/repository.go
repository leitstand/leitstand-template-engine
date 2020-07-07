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

package job

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/rs/zerolog/log"
)

//go:generate moq -out RepositoryMock.go  . Repository
//Repository ...
type Repository interface {
	//Jobs returns a list of jobs
	Jobs() (map[string]*Job, error)
	//Job returns a job or nil
	Job(id string) (*Job, error)
	//AddJob adds a new Job
	AddJob(job *Job) error
	//MakeCallbackToURI Initiate the callback
	MakeCallbackToURI(responseURI string, job *Job)
	//WriteJobResult write interface as data
	WriteJobResult(w http.ResponseWriter, statusCode int, job *Job)
}

//DefaultRepository ...
type DefaultRepository struct {
	jobs            *TTLMap
	httpRetryClient *retryablehttp.Client
	restBaseURL     string
}

//NewDefaultRepository ...
func NewDefaultRepository(restBaseURL string) (r Repository) {
	httpRetryClient := retryablehttp.NewClient()

	return &DefaultRepository{

		jobs:            NewTTLMap(20, 5*time.Minute),
		httpRetryClient: httpRetryClient,
		restBaseURL:     restBaseURL,
	}

}

//Jobs returns a list of jobs
func (m *DefaultRepository) Jobs() (map[string]*Job, error) {
	theMap := m.jobs.Map()
	result := make(map[string]*Job, len(theMap))
	for k, v := range m.jobs.Map() {
		result[k] = v.(*Job)
	}
	return result, nil
}

//Job returns a job or nil
func (m *DefaultRepository) Job(id string) (*Job, error) {
	job, ok := m.jobs.Get(id)
	if !ok {
		return nil, nil
	}
	return job.(*Job), nil
}

//AddJob adds a new Job
func (m *DefaultRepository) AddJob(job *Job) error {
	m.jobs.Put(job.ID, job)
	return nil
}

//MakeCallbackToURI Initiate the callback
func (m *DefaultRepository) MakeCallbackToURI(responseURI string, job *Job) {
	if responseURI != "" {
		var writer bytes.Buffer
		jsonEncoder := json.NewEncoder(&writer)
		_ = jsonEncoder.Encode(job)

		//Create a request
		req, err := retryablehttp.NewRequest(http.MethodPut, responseURI, &writer)
		if err != nil {
			log.Error().Err(err).Msg("")
			return
		}
		req.Header.Set("Content-Type", "application/json")
		res, err := m.httpRetryClient.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer func() { _ = res.Body.Close() }()
		if res.StatusCode < 200 || res.StatusCode >= 300 {
			log.Info().Int("status", res.StatusCode).Msg("not able successfully to do the job callback")
		}
	}
}

//WriteJobResult write interface as data
func (m *DefaultRepository) WriteJobResult(w http.ResponseWriter, statusCode int, job *Job) {
	WriteJobResult(w, statusCode, job, m.restBaseURL)
}

//WriteJobResult write interface as data
func WriteJobResult(w http.ResponseWriter, statusCode int, job *Job, baseURL string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("%s/%s", baseURL, job.ID))
	w.WriteHeader(statusCode)
	jsonEncoder := json.NewEncoder(w)
	_ = jsonEncoder.Encode(job)
}
