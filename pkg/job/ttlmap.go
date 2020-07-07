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
	"sync"
	"time"
)

type ttlItem struct {
	value      interface{}
	lastAccess int64
}

//TTLMap deletes the items after a certain time after insertion
type TTLMap struct {
	m     map[string]*ttlItem
	mutex sync.RWMutex
}

//NewTTLMap will create an new TTLMap. Where ln is the initial capacity and maxTTL is the duration in seconds how long the entries should survive.
func NewTTLMap(ln int, maxTTL time.Duration) (m *TTLMap) {
	maxTTLInSeconds := int64(maxTTL.Seconds())
	m = &TTLMap{m: make(map[string]*ttlItem, ln)}
	go func() {
		for now := range time.Tick(maxTTL / 10) {
			expireReference := now.Unix()
			m.mutex.Lock()
			for k, v := range m.m {
				if v.lastAccess+maxTTLInSeconds < expireReference {
					delete(m.m, k)
				}
			}
			m.mutex.Unlock()
		}
	}()
	return
}

//Len of the map
func (m *TTLMap) Len() int {
	return len(m.m)
}

//Put value to the map
func (m *TTLMap) Put(k string, v interface{}) {
	m.mutex.Lock()
	_, ok := m.m[k]
	if !ok {
		it := &ttlItem{value: v, lastAccess: time.Now().Unix()}
		m.m[k] = it
	}
	m.mutex.Unlock()
}

//Get the value of the map for the key
func (m *TTLMap) Get(k string) (interface{}, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if it, ok := m.m[k]; ok {
		it.lastAccess = time.Now().Unix()
		return it.value, ok

	}
	return nil, false
}

//Map returns the whole map
func (m *TTLMap) Map() map[string]interface{} {
	m.mutex.RLock()
	result := make(map[string]interface{}, len(m.m))
	for key, value := range m.m {
		result[key] = value.value
	}
	m.mutex.RUnlock()
	return result
}
