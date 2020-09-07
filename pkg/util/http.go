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

package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// WriteMessage writes the particular message to the response
func WriteMessage(w http.ResponseWriter, statuscode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statuscode)
	jsonEncoder := json.NewEncoder(w)
	_ = jsonEncoder.Encode(&Message{Message: message})
}

// WriteAsJSON write interface as data
func WriteAsJSON(w http.ResponseWriter, statuscode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statuscode)
	jsonEncoder := json.NewEncoder(w)
	_ = jsonEncoder.Encode(data)
}

// ReadJSON write interface as data
func ReadJSON(req *http.Request, data interface{}) error {
	decoder := json.NewDecoder(req.Body)
	return decoder.Decode(data)
}

// ReadJSONObject reads the file under the filename as JSON object
func ReadJSONObject(filename string, object interface{}) error {
	// Open our jsonFile
	jsonFile, err := os.Open(filename)
	// if we os.Open returns an error then handle it
	if err != nil {
		return err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer func() { _ = jsonFile.Close() }()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	// var result BdsObject
	return json.Unmarshal([]byte(byteValue), object)
}

//ValidateAndGetVariableFromPath ...
func ValidateAndGetVariableFromPath(w http.ResponseWriter, req *http.Request, variableName string) (string, bool) {
	vars := mux.Vars(req)
	name, set := vars[variableName]
	if !set {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error: %s: There is an misconfiguration in the path variables!\n", variableName)
		return name, false
	}
	return name, true
}
