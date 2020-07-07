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

// GenerationRequest to generate a config
type GenerationRequest struct {
	//PutBackURL where the result should be sent back, only used for the async call
	PutBackURL string `json:"put_back_url"`
	//Variables for the generation
	Variables map[string]interface{} `json:"variables"`
}
