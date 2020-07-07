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

package options

import (
	"fmt"
	"strings"
	"testing"

	isTest "github.com/matryer/is"
)

func testOptions() *Options {
	o := &Options{
		HTTPAddress: "http://127.0.0.1:29092",
	}
	return o
}

func errorMsg(msgs []string) string {
	err := fmt.Errorf("%w\ndetail:\n%s", ErrInvalidConfiguration, strings.Join(msgs, "\n  "))
	return err.Error()
}

func TestNoError(t *testing.T) {
	is := isTest.New(t)
	o := &Options{
		HTTPAddress:  "http://localhost:29092",
		TemplatePath: "./templates",
	}
	err := o.Validate()
	is.NoErr(err)
}

func TestTemplatePathOptions(t *testing.T) {
	expected := errorMsg([]string{
		"missing setting: template_path",
	})
	is := isTest.New(t)
	o := testOptions()
	err := o.Validate()
	is.True(err != nil)
	if err != nil {
		is.Equal(expected, err.Error())
	}
}

func TestNewOptions(t *testing.T) {
	expected := errorMsg([]string{
		"missing setting: http-address",
		"missing setting: template_path",
	})
	is := isTest.New(t)
	o := &Options{}
	err := o.Validate()
	is.True(err != nil)
	if err != nil {
		is.Equal(expected, err.Error())
	}
}
