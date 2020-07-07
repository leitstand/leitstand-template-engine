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

package configen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/rs/zerolog/log"

	"github.com/leitstand/leitstand-template-engine/pkg/configen"
	"github.com/leitstand/leitstand-template-engine/pkg/util"

	isTest "github.com/matryer/is"

	"github.com/google/go-cmp/cmp"
	"github.com/yosuke-furukawa/json5/encoding/json5"
)

// TestRepository_IntegrationTest are all tests for the sample templates folder.
// This ensures code changes does not affect the samples.
func TestRepository_IntegrationTest(t *testing.T) {
	tests := []struct {
		name         string
		templatePath string
		template     string
		test         string
		format       string
	}{
		{
			templatePath: "../../templates",
			template:     "sample",
			test:         "example",
			format:       "json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := isTest.New(t)
			r := configen.NewRepository(tt.templatePath)

			var variables map[string]interface{}
			variablesFile := fmt.Sprintf("%s/%s/%s_variables.%s", tt.templatePath, tt.template, tt.test, tt.format)
			err := util.ReadJSONObject(variablesFile, &variables)
			if err != nil {
				log.Error().Err(err).Msg("can't find variables file")
				return
			}
			resultFile := fmt.Sprintf("%s/%s/%s_result.%s", tt.templatePath, tt.template, tt.test, tt.format)
			want, err := ioutil.ReadFile(resultFile)
			is.NoErr(err)

			got, err := r.GenerateFile(tt.template, variables)
			if err != nil {
				log.Error().Err(err).Msg("error in generation")
				return
			}
			gotFile := fmt.Sprintf("%s/%s/%s_got.%s", tt.templatePath, tt.template, tt.test, tt.format)
			_ = ioutil.WriteFile(gotFile, got, os.ModePerm)
			switch tt.format {
			case "txt":
				if diff := cmp.Diff(string(want), string(got)); diff != "" {
					t.Errorf("mismatch (-want +got):\n%s", diff)
				}
			case "json":
				if diff := cmp.Diff(transformToJSONObject(want), transformToJSONObject(got)); diff != "" {
					t.Errorf("mismatch (-want +got):\n%s", diff)
				}
			case "json5":
				if diff := cmp.Diff(transformToJSON5Object(want), transformToJSON5Object(got)); diff != "" {
					t.Errorf("mismatch (-want +got):\n%s", diff)
				}
			default:
				t.Error("unknown format")
			}
		})
	}
}
func transformToJSONObject(value []byte) interface{} {
	var v interface{}
	if err := json.Unmarshal(value, &v); err != nil {
		return fmt.Sprintf("not parseable (%s)", value) // use unparseable input as the output
	}
	return v
}
func transformToJSON5Object(value []byte) interface{} {
	var v interface{}
	if err := json5.Unmarshal(value, &v); err != nil {
		return value // use unparseable input as the output
	}
	return v
}
