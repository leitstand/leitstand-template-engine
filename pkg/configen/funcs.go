/*
 * Copyright 2022 RtBrick Inc.
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
	"fmt"
	"reflect"

	sv2 "github.com/Masterminds/semver"
)

func featureIsEnabled(constraint, version interface{}) (bool, error) {

	// If no feature toggle is specied the feature is disable by default.
	if constraint == nil {
		return false, nil
	}

	// The constraint is supposed to be a string.
	s, ok := constraint.(string)
	if !ok {
		return false, fmt.Errorf("invalid feature toggle constraint type: %v", reflect.TypeOf(constraint))
	}

	// An empty feature toggle value is handled like a nil value.
	if s == "" {
		return false, nil
	}

	// Parse the given constraint and report invalid constraints to the template engine.
	c, err := sv2.NewConstraint(s)
	if err != nil {
		return false, err
	}

	// Parse version and report invalid version value to the template engine.
	v, err := sv2.NewVersion(version.(string))
	if err != nil {
		return false, err
	}

	// Check whether the version matches the expected range.
	return c.Check(v), nil
}
