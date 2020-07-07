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
	"fmt"
	"log"
	"net/http"

	"github.com/leitstand/leitstand-template-engine/pkg/job"
	"github.com/leitstand/leitstand-template-engine/pkg/util"

	"github.com/hashicorp/go-retryablehttp"
)

// @Summary generate a configuration file
// @Description generate a configuration file
// @Description **Characteristics:**
// @Description * Operation: **asynchronous**
// @Tags template-engine
// @Accept  json
// @Produce  json
// @Param response_uri header string false "callback response uri"
// @Param template_name path string true "name of the template"
// @Param body body GenerationRequest true "body"
// @Header 202 {string} Location "Location to get the job result"
// @Success 202 "Accepted"
// @Failure 422 {object} util.Message
// @Failure 500 {object} util.Message
// @Router /template-engine/api/v1/templates/{template_name}/_generate [POST]
func (app *Application) generateConfigurationAsync(w http.ResponseWriter, req *http.Request) {
	responseURI := req.Header.Get("response_uri")
	templateName, ok := util.ValidateAndGetVariableFromPath(w, req, "template_name")
	if !ok {
		return
	}
	requestBody := &GenerationRequest{}
	err := util.ReadJSON(req, requestBody)
	if err != nil {
		log.Printf("Error: %s\n", err)
		util.WriteMessage(w, http.StatusBadRequest, fmt.Sprintf("error %v", err))
		return
	}

	asyncJob := job.NewJob(fmt.Sprintf("generate configuration: %s", templateName))
	_ = app.jobRepository.AddJob(asyncJob)
	app.jobRepository.WriteJobResult(w, http.StatusAccepted, asyncJob)
	go func() {
		defer app.jobRepository.MakeCallbackToURI(responseURI, asyncJob)
		result, err := app.repository.GenerateFile(templateName, requestBody.Variables)
		if err != nil {
			asyncJob.SetResult(job.NewAsyncResultWithMessage(http.StatusBadRequest, fmt.Sprintf("error %v", err)))
			return
		}

		err = makeCallbackToURI(requestBody.PutBackURL, result)
		if err != nil {
			asyncJob.SetResult(job.NewAsyncResultWithMessage(http.StatusBadRequest, fmt.Sprintf("error %v", err)))
			return
		}
		asyncJob.SetResult(job.NewAsyncResult(http.StatusOK))
	}()

}

//makeCallbackToURI Initiate the callback
func makeCallbackToURI(responseURI string, data []byte) error {
	if responseURI != "" {
		// Create a request
		req, err := retryablehttp.NewRequest("PUT", responseURI, data)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return nil
		}

		req.Header.Set("Content-Type", "application/json")
		resp, err := retryablehttp.NewClient().Do(req)
		if err == nil {
			defer func() {
				_ = resp.Body.Close()
			}()
		} else {
			return err
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("not able successfully to do the callback. Status code %d", resp.StatusCode)
		}
	}
	return nil
}

// @Summary generate a configuration file
// @Description generate a configuration file
// @Description **Characteristics:**
// @Description * Operation: **synchronous**
// @Tags template-engine
// @Accept  json
// @Produce  json
// @Param template_name path string true "name of the template"
// @Param body body GenerationRequest true "body"
// @Success 200 "config file"
// @Failure 422 {object} util.Message
// @Failure 500 {object} util.Message
// @Router /template-engine/api/v1/templates/{template_name}/_generatesync [POST]
func (app *Application) generateConfigurationSync(w http.ResponseWriter, req *http.Request) {
	templateName, ok := util.ValidateAndGetVariableFromPath(w, req, "template_name")
	if !ok {
		return
	}
	requestBody := &GenerationRequest{}
	err := util.ReadJSON(req, requestBody)
	if err != nil {
		log.Printf("Error: %s\n", err)
		util.WriteMessage(w, http.StatusBadRequest, fmt.Sprintf("error %v", err))
		return
	}

	result, err := app.repository.GenerateFile(templateName, requestBody.Variables)
	if err != nil {
		util.WriteMessage(w, http.StatusBadRequest, fmt.Sprintf("error %v", err))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(result)
}
