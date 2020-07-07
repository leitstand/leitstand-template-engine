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
	"math/rand"
	"strconv"

	"github.com/leitstand/leitstand-template-engine/pkg/util"
)

type JobState string

const (
	//StatusPending job status
	StatusPending JobState = "PENDING"
	//StatusDone job status
	StatusDone JobState = "DONE"
)

//Job ...
type Job struct {
	ID          string   `json:"id"`                         //Id of the job
	State       JobState `json:"state" enums:"PENDING,DONE"` //State of the Job.
	Description string   `json:"description,omitempty"`      //Description of the Job.
	Result      *Result  `json:"result,omitempty"`
	//CallbackURL allows to store the callback url if the callback has to be done later
	CallbackURL string `json:"-"`
}

//SetResult sets the result of the job
func (j *Job) SetResult(result *Result) {
	j.Result = result
	j.State = StatusDone
}

//Result ...
type Result struct {
	Status int         `json:"status"`         //Http Code if it would be in an synchronous call
	Data   interface{} `json:"data,omitempty"` //Body data if it would be in an synchronous call
}

//NewAsyncResultWithMessage creates the particular message as Result
func NewAsyncResultWithMessage(statusCode int, message string) *Result {
	return &Result{
		Status: statusCode,
		Data:   &util.Message{Message: message},
	}
}

//NewAsyncResult creates the particular message as Result
func NewAsyncResult(statusCode int) *Result {
	return &Result{
		Status: statusCode,
	}
}

//NewJob ...
func NewJob(description string) *Job {
	return &Job{
		ID:          strconv.Itoa(rand.Int()),
		Description: description,
		State:       StatusPending,
	}
}
