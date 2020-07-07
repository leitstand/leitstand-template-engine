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
	"errors"
	"fmt"
	"strings"
)

var (
	//ErrInvalidConfiguration invalid configuration
	ErrInvalidConfiguration = errors.New("invalid configuration")
)

// Options for the leitstand-template-engine
type Options struct {
	HTTPAddress  string `json:"http_address"`
	TemplatePath string `json:"template_path"`
}

// Validate the options
func (o *Options) Validate() error {
	msgs := make([]string, 0)
	if len(o.HTTPAddress) < 1 {
		msgs = append(msgs, "missing setting: http-address")
	}
	// if only one is set
	if len(o.TemplatePath) < 1 {
		msgs = append(msgs, "missing setting: template_path")
	}

	if len(msgs) != 0 {
		return fmt.Errorf("%w\ndetail:\n%s", ErrInvalidConfiguration,
			strings.Join(msgs, "\n  "))
	}
	return nil
}
