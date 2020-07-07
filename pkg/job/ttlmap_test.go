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
	"fmt"
	"testing"
	"time"
)

func TestTTLMap(t *testing.T) {
	ttlmap := NewTTLMap(10, time.Second)
	if ttlmap.Len() != 0 {
		t.Error("initial len should be 0")
	}
	ttlmap.Put("key1", "value1")
	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Microsecond)
		value, ok := ttlmap.Get("key1")
		if !ok || value != "value1" {
			t.Error("element not found anymore")
		}
	}
	time.Sleep(2 * time.Second)
	v, ok := ttlmap.Get("key1")
	fmt.Println(v, ok)
	if ok {
		t.Error("element should not be found anymore")
	}
}

func TestTTLMap_Map(t *testing.T) {
	ttlmap := NewTTLMap(10, time.Minute)
	if ttlmap.Len() != 0 {
		t.Error("initial len should be 0")
	}
	ttlmap.Put("key1", "value1")
	realmap := ttlmap.Map()
	if len(realmap) != ttlmap.Len() {
		t.Error("maps should have same length")
	}
	for key, value := range realmap {
		ttl, _ := ttlmap.Get(key)
		if ttl != value {
			t.Error("Values of maps should be the same")
		}
	}
}
func TestTTLMap_Put(t *testing.T) {
	ttlmap := NewTTLMap(10, time.Second)
	if ttlmap.Len() != 0 {
		t.Error("initial len should be 0")
	}
	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Microsecond)
		ttlmap.Put("key1", "value1")
	}
	value, ok := ttlmap.Get("key1")
	if !ok || value != "value1" {
		t.Error("element not found anymore")
	}
	time.Sleep(2 * time.Second)
	v, ok := ttlmap.Get("key1")
	fmt.Println(v, ok)
	if ok {
		t.Error("element should not be found anymore")
	}
}
