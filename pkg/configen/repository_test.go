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
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/leitstand/leitstand-template-engine/pkg/util"

	isTest "github.com/matryer/is"

	"github.com/google/go-cmp/cmp"
)

func Test_parseConfigFile(t *testing.T) {
	type args struct {
		templatePath   string
		templateFolder string
	}
	tests := []struct {
		name      string
		args      args
		want      *TemplateConfig
		wantErr   bool
		wantedErr error
	}{
		{
			args:      args{templatePath: "testdata/templates", templateFolder: "t0"},
			wantErr:   true,
			wantedErr: ErrTemplateConfigNotFound,
		},
		{
			args: args{templatePath: "testdata/templates", templateFolder: "t1"},
			want: &TemplateConfig{
				TemplateEngine: "golang",
				MainPattern:    "testdata/templates/t1/*.goyaml",
				IncludePattern: "testdata/templates/includes/*.goyaml",
				MainTemplate:   "main.goyaml",
			},
		}, {
			args: args{templatePath: "testdata/templates", templateFolder: "t2"},
			want: &TemplateConfig{
				TemplateEngine: "golang",
				MainPattern:    "",
				IncludePattern: "",
				MainTemplate:   "main.goyaml",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := isTest.New(t)
			got, err := parseConfigFile(tt.args.templatePath, tt.args.templateFolder)
			if tt.wantErr {
				if err == nil {
					t.Errorf("parseConfigFile() error = %v, wantErr %v", err, tt.wantErr)
				}
				if tt.wantedErr != nil && !errors.Is(err, tt.wantedErr) {
					t.Errorf("parseConfigFile() error = %v, wantedErr %v", err, tt.wantedErr)
				}
				return
			}
			is.NoErr(err)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestRepository_GenerateFile(t *testing.T) {
	tests := []struct {
		name           string
		templatePath   string
		templateFolder string
		want           []byte
		wantErr        bool
		wantedErr      error
	}{
		{
			templatePath:   "testdata/templates",
			templateFolder: "g0",
			wantErr:        true,
			wantedErr:      ErrTemplateConfigNotFound,
		}, {
			templatePath:   "testdata/templates",
			templateFolder: "g1",
			wantErr:        true,
			wantedErr:      ErrEngineNotFound,
		}, {
			templatePath:   "testdata/templates",
			templateFolder: "g2",
			wantErr:        false,
			want:           []byte("Hi Chris!\nfooter"),
		}, {
			templatePath:   "testdata/templates",
			templateFolder: "g3",
			wantErr:        false,
			want:           []byte(`{"name":"Chris","static":"static"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := isTest.New(t)
			r := NewRepository(tt.templatePath)

			var variables map[string]interface{}
			variablesFile := fmt.Sprintf("%s/%s/%s", tt.templatePath, tt.templateFolder, "variables.json")
			_ = util.ReadJSONObject(variablesFile, &variables)

			got, err := r.GenerateFile(tt.templateFolder, variables)
			if tt.wantErr {
				if err == nil {
					t.Errorf("parseConfigFile() error = %v, wantErr %v", err, tt.wantErr)
				}
				if tt.wantedErr != nil && !errors.Is(err, tt.wantedErr) {
					t.Errorf("parseConfigFile() error = %v, wantedErr %v", err, tt.wantedErr)
				}
				return
			}
			is.NoErr(err)
			if diff := cmp.Diff(string(tt.want), string(got)); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_removeTrailingCommas(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    []byte
		wantErr bool
	}{
		{
			input:   []byte(`{"key1":"v1","key2":"v2",}`),
			want:    []byte(`{"key1":"v1","key2":"v2"}`),
			wantErr: false,
		},
		{
			input:   []byte(`{"key1":"v1","key2":"v2"  ,  }`),
			want:    []byte(`{"key1":"v1","key2":"v2"    }`),
			wantErr: false,
		}, {
			input: []byte(`{"key1":"v1","key2":"v2"
,
}`),
			want: []byte(`{"key1":"v1","key2":"v2"

}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := removeTrailingCommas(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("removeTrailingCommas() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeTrailingCommas() got = %s, want %s", got, tt.want)
			}
		})
	}
}

func Test_removeEmptyLines(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    []byte
		wantErr bool
	}{
		{
			input: []byte(`{
 	"key1":"v1",
	"key2":"v2"
}`),
			want: []byte(`{
 	"key1":"v1",
	"key2":"v2"
}`),
			wantErr: false,
		},
		{
			input: []byte(`{
 	"key1":"v1",


	"key2":"v2"
  
  


}`),
			want: []byte(`{
 	"key1":"v1",
	"key2":"v2"
}`),
			wantErr: false,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := removeEmptyLines(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("removeTrailingCommas() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeTrailingCommas() got = %s, want %s", got, tt.want)
			}
		})
	}
}

func Test_prettyJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    []byte
		wantErr bool
	}{
		{
			input:   []byte(`{"key1": "v1","key2": "v2"}`),
			want:    []byte("{\n  \"key1\": \"v1\",\n  \"key2\": \"v2\"\n}\n"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := prettyJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("removeTrailingCommas() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(string(tt.want), string(got)); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}
func Test_uglyJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    []byte
		wantErr bool
	}{
		{
			input:   []byte("{\n  \"key1\": \"v1\",\n  \"key2\": \"v2\"\n}\n"),
			want:    []byte(`{"key1":"v1","key2":"v2"}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := uglyJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("removeTrailingCommas() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(string(tt.want), string(got)); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}
